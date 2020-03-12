# go-iabtcf
A Golang implementation of the IAB Transparency and Consent String (TC String) v2.0

The TC String is a technical component of the IAB Europe Transparency & Consent Framework (TCF)

Library provides convenient way to check if:
- vendor is allowed
- purposes are allowed
- special fetures are allowed

## Installation

    go get -v github.com/travelaudience/go-iabtcf
    
## Example

    package main
    
    import (
      "fmt"
    
      "github.com/travelaudience/go-iabtcf"
    )
    
    func main() {
      var s, err = iabtcf.ParseCoreString("COwIsAvOwIsAvBIAAAENAPCMAP_AAP_AAAAAFoQBQABAAGAAQAAwACQAAAAA.IFoEUQQgAIQwgIwQABAEAAAAOIAACAIAAAAQAIAgEAACEAAAAAgAQBAAAAAAAGBAAgAAAAAAAFAAECAAAgAAQARAEQAAAAAJAAIAAgAAAYQEAAAQmAgBC3ZAYzUw")
      if err != nil {
        panic(err)
      }
      
      pa := s.EveryPurposeAllowed([]int{1})
      sf := s.EverySpecialFeatureAllowed([]int{1})
      va := s.VendorAllowed(1)
    }