# goup
go practises

reusable go functions

examples as follows:

  Any([]int{1,2,3}, func(n int)bool{ return n%2==1})

  Find([]int{1,2,3}, func(n int)bool{ return n%2==1})

  Filter([]int{1,2,3}, func(n int)bool{ return n%2==1})

  Map([]int{1,2,3}, func(n int)int{ return n*n })

  Reduce([]int{1,2,3}, func(a,b int)int{ return a+b })

  for idx := range CRange(9,1,2){ 
    fmt.Printf("%d\n", idx) 
  }// 1,3,5,7 




