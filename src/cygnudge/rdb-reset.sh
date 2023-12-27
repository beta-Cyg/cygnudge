#!/bin/sh

redis-cli flushdb
redis-cli set user-number 0
