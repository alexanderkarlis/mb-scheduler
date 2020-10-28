# Mindbody Scheduler

#### Mindbody-scheduler is a scheduler to run every alotted time to sign up for classes on [mindbody.io](mindbody.io)

## Requirements
Uses geckodriver (firefox) as the webdriver
### (Unix)
* Go 1.14+ 
* Xfvb
* Java SDK

### (Windows)
* GO 1.14+ 
* Java SDK

## Run
```sh
> export MINDBODY_USERNAME="username@test.com"
> export MINDBODY_PASSWORD="secretpassword"
> 
> go build
>
> ./mindbody -date="10/15/2020" -time="5:45pm" -fullname="Alexander Karlis" \
-username=$MINDBODY_USERNAME -password=$MINDBODY_PASSWORD
> # on unix
> eog sc  # to check the screenshot at the end for confirmation
