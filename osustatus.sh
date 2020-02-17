#!/bin/bash


scanmem -p `pgrep osu\\!.exe` -e -c 'option scan_data_type bytearray;48 83 F8 04 73 1E;list;exit'
