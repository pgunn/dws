#!/bin/bash

dropdb dws
createdb dws
psql dws < mk_db.sql
psql dws < db_load_sampledata.sql
if [ -f ~/local_dws.sql ]; then
	psql dws < ~/local_dws.sql
fi
