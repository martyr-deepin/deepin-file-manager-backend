/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#include <gio/gio.h>

int content_type_can_be_executable(char* type)
{
    return g_content_type_can_be_executable(type);
}


int content_type_is(char* type, char* expected_type)
{
    return g_content_type_is_a(type, expected_type);
}

