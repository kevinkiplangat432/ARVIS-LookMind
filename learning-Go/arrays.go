package main

import ("fmt")

func arrayDemo() {
  var arr1 = [3]int{1,2,3} // length is defined 
  arr2 := [5]int{4,5,6,7,8}
  var arr3 = [...]string{"kevin", "laban", "alvin", "angella", "mr who"}  // the lenght is defined
 
  fmt.Println(arr1)
  fmt.Println(arr2)
  fmt.Println(arr3)
}

func array_of_strings(){
	var cars = [5]string{"Bmw", "volvo", "volkswagon","benz", "lambo"}
	fmt.Println(cars)
}


// accessing array elements
func accessing_elements(){
	prices := [3]int{10,20,30}

	fmt.Println(prices[0])
	fmt.Println(prices[2])


}


func change_array_elements(){
	price := [3]int{10,20,30}
	price[2] = 50
	fmt.Println(price)
}


func array_initialisation() {
  arr1 := [5]int{} //not initialized
  arr2 := [5]int{1,2} //partially initialized
  arr3 := [5]int{1,2,3,4,5} //fully initialized
  arr4 := [5]int{1:10,2:40}  // initialize specific ELements like the seconsd and the third

  fmt.Println(arr1)
  fmt.Println(arr2)
  fmt.Println(arr3)
  fmt.Println(arr4)
  fmt.Println(len(arr1))
}