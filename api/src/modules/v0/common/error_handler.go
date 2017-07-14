//Author: Peter Nagy <https://peternagy.ie>
//Since: 06, 2017
//Description: handles errors
package common

import "log"

//CheckError - Generalized check for errors
func CheckError(src string, e error, isPanic bool) {
	if e != nil {
		log.Println(src, e)
		if isPanic {
			panic(e.Error)
		}
		return
	}
}
