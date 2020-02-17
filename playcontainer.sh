#!/bin/bash


scanmem -p `pgrep osu\\!.exe` -e -c 'option scan_data_type bytearray;85 C9 74 1F 8D 55 F0 8B 01;list;exit'
