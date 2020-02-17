#!/bin/bash


scanmem -p `pgrep osu\\!.exe` -e -c 'option scan_data_type bytearray;F8 01 74 04 83;list;exit'
