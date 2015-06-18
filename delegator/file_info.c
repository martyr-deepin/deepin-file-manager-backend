#include <string.h>
#include <gtk/gtk.h>
#include <gio/gdesktopappinfo.h>

#define APP_ICON_SIZE 48

#define is_xpm(ext) (g_ascii_strcasecmp(ext, "xpm") == 0)
#define is_dataurl(ext) g_str_has_prefix(ext, "data:image")
#define is_deepin_icon(icon_path) g_str_has_prefix(icon_path, "/usr/share/icons/Deepin/")


static
char* get_data_uri_by_pixbuf(GdkPixbuf* pixbuf)
{
    gchar* buf = NULL;
    gsize size = 0;
    GError *error = NULL;

    gdk_pixbuf_save_to_buffer(pixbuf, &buf, &size, "png", &error, NULL);
    g_assert(buf != NULL);

    if (error != NULL) {
        g_warning("[%s] %s\n", __func__, error->message);
        g_error_free(error);
        g_free(buf);
        return NULL;
    }

    char* base64 = g_base64_encode((const guchar*)buf, size);
    g_free(buf);
    char* data = g_strconcat("data:image/png;base64,", base64, NULL);
    g_free(base64);

    return data;
}


static
char* get_data_uri_by_path(const char* path)
{
    GError *error = NULL;
    GdkPixbuf* pixbuf = gdk_pixbuf_new_from_file(path, &error);
    if (error != NULL) {
        g_warning("%s\n", error->message);
        g_error_free(error);
        return NULL;
    }
    char* c = get_data_uri_by_pixbuf(pixbuf);
    g_object_unref(pixbuf);
    return c;

}


static
char* icon_name_to_path(const char* name, int size)
{
    if (g_path_is_absolute(name))
        return g_strdup(name);

    g_return_val_if_fail(name != NULL, NULL);

    int pic_name_len = strlen(name);
    char* ext = strrchr(name, '.');
    if (ext != NULL) {
        if (g_ascii_strcasecmp(ext+1, "png") == 0 || g_ascii_strcasecmp(ext+1, "svg") == 0 || g_ascii_strcasecmp(ext+1, "jpg") == 0) {
            pic_name_len = ext - name;
            g_debug("desktop's Icon name should an absoulte path or an basename without extension");
        }
    }

    char* pic_name = g_strndup(name, pic_name_len);
    GtkIconTheme* them = gtk_icon_theme_get_default(); // NB: do not ref or unref it

    GtkIconInfo* info = gtk_icon_theme_lookup_icon(them, pic_name, size, GTK_ICON_LOOKUP_GENERIC_FALLBACK);
    g_free(pic_name);

    if (info) {
        char* path = g_strdup(gtk_icon_info_get_filename(info));

#if GTK_MAJOR_VERSION >= 3
        g_object_unref(info);
#elif GTK_MAJOR_VERSION == 2
        gtk_icon_info_free(info);
#endif
        g_debug("get icon from icon theme is: %s", path);
        return path;
    }
    g_warning("get gtk icon theme info failed");

    return NULL;
}


static
char* check_xpm(const char* path)
{
    if (path == NULL)
        return NULL;
    char* ext = strrchr(path, '.');
    if (ext != NULL && is_xpm(ext+1)) {
        return get_data_uri_by_path(path);
    } else {
        return g_strdup(path);
    }
}


static
char* icon_name_to_path_with_check_xpm(const char* name, int size)
{
    char* path = icon_name_to_path(name, size);
    char* icon = check_xpm(path);
    g_free(path);
    return icon;
}


static
char* get_basename_without_extend_name(char const* path)
{
    g_assert(path!= NULL);
    char* basename = g_path_get_basename(path);
    char* ext_sep = strrchr(basename, '.');
    if (ext_sep != NULL) {
        char* basename_without_ext = g_strndup(basename, ext_sep - basename);
        g_free(basename);
        return basename_without_ext;
    }

    return basename;
}


static
char* _check(char const* app_id)
{
    char* icon = NULL;
    char* temp_icon_name_holder = icon_name_to_path_with_check_xpm(app_id, 48);

    if (temp_icon_name_holder != NULL) {
        if (!is_dataurl(temp_icon_name_holder))
            icon = temp_icon_name_holder;
        else
            g_free(temp_icon_name_holder);
    }

    return icon;
}


static
char* check_absolute_path_icon(char const* app_id, char const* icon_path)
{
    char* icon = NULL;
    if ((icon = _check(app_id)) == NULL) {
        char* basename = get_basename_without_extend_name(icon_path);
        if (basename != NULL) {
            if (g_strcmp0(app_id, basename) == 0 || (icon = _check(basename)) == NULL) {
                icon = g_strdup(icon_path);
            }
            g_free(basename);
        }
    }

    return icon;
}


char* get_icon_from_app(char const* file_path)
{
    GDesktopAppInfo* app = g_desktop_app_info_new_from_filename(file_path);
    if (app == NULL) {
        return NULL;
    }

    // NB: g_app_info_get_icon transfer none, do NOT unref icon.
    GIcon* gicon = g_app_info_get_icon(G_APP_INFO(app));
    if (gicon == NULL) {
        return NULL;
    }

    char* icon_str = g_icon_to_string(gicon);
    g_debug("app icon: %s", icon_str);
    if (icon_str != NULL && g_path_is_absolute(icon_str) && !is_deepin_icon(icon_str)) {
        g_debug("check_absolute_path_icon");
        char* app_id = get_basename_without_extend_name(g_desktop_app_info_get_filename(G_DESKTOP_APP_INFO(app)));
        char* temp_icon_name_holder = icon_str;
        icon_str = check_absolute_path_icon(app_id, temp_icon_name_holder);
        g_free(app_id);
        g_free(temp_icon_name_holder);
    }

    char* icon = icon_name_to_path_with_check_xpm(icon_str, APP_ICON_SIZE);
    g_free(icon_str);
    g_debug("the final icon of app is: %s", icon);

    g_object_unref(app);
    return icon;
}


char* get_icon_from_file(char* icons)
{
    if (icons == NULL) {
        return NULL;
    }

    char* icon = NULL;
    char** icon_names = g_strsplit(icons, " ", -1);

    for (int i = 0; icon_names[i] != NULL && icon == NULL; ++i) {
        icon = icon_name_to_path(icon_names[i], APP_ICON_SIZE);
    }

    g_strfreev(icon_names);

    return icon;
}

