package main;




func main(){
	port := os.Getenv("PORT");

	if port == ""{
		port = "8080";
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port);
	}

	
}