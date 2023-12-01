package main

import "strings"

// replace with frontend app URL
const OriginURL = "http://localhost:3000"

// replace with your table name
const TableName = "Bananas"

var ApiPath = strings.ToLower(TableName)
