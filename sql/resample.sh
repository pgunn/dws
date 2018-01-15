#!/bin/bash

dropdb dws
createdb dws
psql dws < mk_db.sql
psql dws < db_load_sampledata.sql
