### Usage
```go
package main

import (
        "fmt"
        "math/big"
        "github.com/abulo/ratel/v1/alpha"
)

func main() {
        input := "123456789012345678901234567890"
        inputBigInt := big.NewInt(0)
        inputBigInt.SetString(input, 10)

        fmt.Println("Input : ", input)

        // encode and decode functions
        encoded, err := alpha.Encode(input)
        if err != nil {
                fmt.Println(err)
        }
        fmt.Println("Encoded : ", encoded)

        decoded, err := alpha.Decode(encoded)
        if err != nil {
                fmt.Println(err)
        }
        fmt.Println("Decoded : ", decoded)

        if input == decoded {
                fmt.Println("Passed! decoded value is the same as the original.")
        } else {
                fmt.Println("FAILED! decoded value is NOT the same as the original!!")
        }

        // encode int and decode int functions
        encodedInt, err := alpha.EncodeInt(inputBigInt)
        if err != nil {
                fmt.Println(err)
        }
        fmt.Println("Encoded using big int: ", encodedInt)

        decodedInt, err := alpha.DecodeInt(encodedInt)
        if err != nil {
                fmt.Println(err)
        }
        fmt.Println("Decoded using big int: ", decodedInt)

        if inputBigInt.Cmp(decodedInt) == 0 {
                fmt.Println("Passed! decoded int value is the same as the intput big int.")
        } else {
                fmt.Println("FAILED! decoded int value is NOT the same as the original!!")
        }
}
```
#### output looks like,

```go
Input :  123456789012345678901234567890

Encoded :  2aYls9bkamJJSwhr0
Decoded :  123456789012345678901234567890
Passed! decoded value is the same as the original.

Encoded using big int:  2aYls9bkamJJSwhr0
Decoded using big int:  123456789012345678901234567890
Passed! decoded int value is the same as the intput big int.
```