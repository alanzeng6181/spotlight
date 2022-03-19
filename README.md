
# Spotlight

A little program that uses QuadTree data structure and fyne ui framework to generate and track dots.
Every 3 seconds, one of the dots is randomly selected, it and its surrounding dots within a specified radius will be highlighted to create spotlight effect.


![output4](https://user-images.githubusercontent.com/17134457/159099300-6a35b2ac-471d-4062-8ca1-64546f743181.gif)


## Usage

run `go build main.go` to build the binary file `main`
  
run `./main -h` to see list of flags:  
  
```
  -c int
        number of dots to populate, between 1 and 500000 (default 10000) \n
  -dps int
        dots per second, between 1 and 10000 (default 500)
  -gc float
        center dot enlargement, between 1.0 and 10.0 (default 5)
  -gd int
        glow duration in miliseconds, bewteen 1000 and 1000 (default 2000)
  -gs float
        surrounding dots enlargement, between 1.0 and 10.0 (default 2)
  -h float
        height of screen, between 100 and 5000 (default 1000)
  -r float
        effective search radius, between 10 and 400 (default 100)
  -w float
        width of screen, between 100 and 5000 (default 1900)
```
