// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * Copyright (C) 2018-2022 SCANOSS.COM
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

// Simple file that demostrate how to call the libhpsm library from C

#include <stdio.h>
#include "libhpsm.h"
/*struct metadata_t{

	int size;
	long int md5;
} ;*/

int main() {
   unsigned char arr[]={230,108,251,147,233,121,242,44,133,94,241,255,76,146,139,25,217,72,95,189,54,12,110,45,136,133,169,40,131,30,101,81,212,67,161,235,78,251,177,123,0,255,76,165,0,9,188,50,0,226,0,51,0,56,0,92,0,237,0,229,153,85,33,0,32,0,171,88,145,92,0,152,0,0,255,76,165,0,9,188,50,0,226,0,51,0,56,0,92,0,237,0,229,153,85,33,0,32,0,171,88,145,92,0,152,0,0,255,0};
   char *md5;
   asprintf(&md5,"04d93973aafdb9e2b3474546378a9085");
 
  struct ranges r = ProcessHPSM(arr,109,md5);

   printf("lines: %s - oss_lines %s", r.local, r.remote);
   
   free(md5);
}
