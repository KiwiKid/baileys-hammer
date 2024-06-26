package main

import (
	"fmt"
)

var baseUrl = "/"
var containerClass = "w-full flex mx-auto bg-gray-200 shadow-xl p-2"
var con = "bg-sepia-200 shadow-xl m-10 rounded-lg"

var pri = " bg-blue-500 p-1 hover:bg-blue-600 text-white font-bold md:p-4 py-1 px-1 rounded-lg hover:scale-105 transition transform ease-out duration-200"

var sec = " text-3xl bg-gray-500 p-1 hover:bg-gray-600 text-white font-bold py-2 px-4 rounded-lg hover:scale-105 transition transform ease-out duration-200"

var add = " text-3xl bg-green-500 hover:bg-green-600 text-white font-bold py-2 px-4 rounded-lg hover:scale-105 transition transform ease-out duration-200"
var del = " text-3xl bg-red-500 hover:bg-red-600 text-white font-bold py-2 px-4 rounded-lg hover:scale-105 transition transform ease-out duration-200"

var bigBtnTxt = "text-3xl py-4 mx-auto items-center justify-center w-full text-center py-2 px-4  rounded-md border hover:bg-opacity-75 focus:outline-none"
var sm = " text-lg mx-auto items-center justify-center w-4/5 py-1 px-2"
var bigPri = fmt.Sprintf("%s %s", bigBtnTxt, pri)
var bigSec = fmt.Sprintf("%s %s", bigBtnTxt, sec)
var bigAdd = fmt.Sprintf("%s %s", bigBtnTxt, add)
var bigDel = fmt.Sprintf("%s %s", bigBtnTxt, del)

var smPri = fmt.Sprintf("%s %s", pri, sm)

var S = fmt.Sprint
var F = fmt.Sprintf
