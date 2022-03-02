package tools

import "log"

// Pretty print a 'header' message to console
func Section(msg string) {
	log.Println("------------------------------")
	log.Println(msg)
	log.Println("------------------------------")
}

// Helper function for checking and logging non-fatal errors
func WriteErr(msg string, err error) {
	if err != nil {
		log.Println(msg, ":", err)
	}
}

// Get type of input
func GetType(x interface{}) {
	log.Printf("Type: %T\n", x)
}
