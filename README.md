Form3 API Client
========================


## Design decisions
Client library was designed to be easily extensible.  Every resource has its
own namespace and new namespaces could be created if needed.  By using
dependency injections I was able to write tests and mock HTTP server.
Interface was design with ease of use in mind and hides complexity from the
user.


## Usage
First you need to initialize new `Client` object with API URL that allows you
to interact with the API:
```go
client := NewClient("http://localhost:8080")
```

Currently there is only one resource implemented - `Organisation/Account`.  It
could be accessed simply by accessing `Client` object property.  Following will
fetch all accounts:
```go
accounts, err := client.Organisation.Accounts.List(nil)
```

Another example would be to create new resource that could be simply achieved
by invoking `Create` method as follow:
```go
newAccount := &Account{
  ID:             "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
  Type:           "accounts",
  OrganisationID: "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
  Attributes:     &Attributes{
    Country:               "GB",
    AccountClassification: "Personal",
  },
}
account, err = client.Organisation.Accounts.Create(newAccount)
```

Its always good idea to check some errors too!
```go
if err != nil {
  fmt.Println(err)
```
`err` is a type of `APIError` contains two fields `StatusCode` and
`ErrorMessage`, both could be accessed by `err.(*APIError).StatusCode`.


## Testing
If you want to run tests, simply execute while being in `form3` directory
```shell
go test
```
