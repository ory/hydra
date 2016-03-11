Reflections
===========

Package reflections provides high level abstractions above the golang reflect library.

Reflect library is very low-level and can be quite complex when it comes to do simple things like accessing a structure field value, a field tag...

The purpose of reflections package is to make developers life easier when it comes to introspect structures at runtime.
Its API is inspired from python language (getattr, setattr, hasattr...) and provides a simplified access to structure fields and tags.

*Reflections is an open source library under the MIT license. Any hackers are welcome to supply ideas, features requests, patches, pull requests and so on: see [Contribute]()*

#### Documentation

Documentation is available at http://godoc.org/github.com/oleiade/reflections


## Installation

#### Into the gopath

```
    go get github.com/oleiade/reflections
```

#### Import it in your code

```go
    import (
        "github.com/oleiade/reflections"
    )
```

## Usage

#### Accessing structure fields

##### GetField

*GetField* returns the content of a structure field. It can be very usefull when
you'd wanna iterate over a struct specific fields values for example. You can whether
provide *GetField* a structure or a pointer to structure as first argument.

```go
    s := MyStruct {
        FirstField: "first value",
        SecondField: 2,
        ThirdField: "third value",
    }

    fieldsToExtract := []string{"FirstField", "ThirdField"}

    for _, fieldName := range fieldsToExtract {
        value, err := reflections.GetField(s, fieldName)
        DoWhatEverWithThatValue(value)
    }
```

##### GetFieldKind

*GetFieldKind* returns the [reflect.Kind](http://golang.org/src/pkg/reflect/type.go?s=6916:6930#L189) of a structure field. It can be used to operate type assertion over a structure fields at runtime.  You can whether provide *GetFieldKind* a structure or a pointer to structure as first argument.

```go
    s := MyStruct{
        FirstField:  "first value",
        SecondField: 2,
        ThirdField:  "third value",
    }

    var firstFieldKind reflect.String
    var secondFieldKind reflect.Int
    var err error

    firstFieldKind, err = GetFieldKind(s, "FirstField")
    if err != nil {
        log.Fatal(err)
    }

    secondFieldKind, err = GetFieldKind(s, "SecondField")
    if err != nil {
        log.Fatal(err)
    }
```

##### GetFieldTag

*GetFieldTag* extracts a specific structure field tag. You can whether provide *GetFieldTag* a structure or a pointer to structure as first argument.


```go
    s := MyStruct{}

    tag, err := reflections.GetFieldTag(s, "FirstField", "matched")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(tag)

    tag, err = reflections.GetFieldTag(s, "ThirdField", "unmatched")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(tag)
```

##### HasField

*HasField* asserts a field exists through structure. You can whether provide *HasField* a structure or a pointer to structure as first argument.


```go
    s := MyStruct {
        FirstField: "first value",
        SecondField: 2,
        ThirdField: "third value",
    }

    // has == true
    has, _ := reflections.HasField(s, "FirstField")

    // has == false
    has, _ := reflections.HasField(s, "FourthField")
```

##### Fields

*Fields* returns the list of a structure field names, so you can access or modify them later on. You can whether provide *Fields* a structure or a pointer to structure as first argument.


```go
    s := MyStruct {
        FirstField: "first value",
        SecondField: 2,
        ThirdField: "third value",
    }

    var fields []string

    // Fields will list every structure exportable fields.
    // Here, it's content would be equal to:
    // []string{"FirstField", "SecondField", "ThirdField"}
    fields, _ = reflections.Fields(s)
```

##### Items

*Items* returns the structure's field name to values map. You can whether provide *Items* a structure or a pointer to structure as first argument.


```go
    s := MyStruct {
        FirstField: "first value",
        SecondField: 2,
        ThirdField: "third value",
    }

    var structItems map[string]interface{}

    // Items will return a field name to
    // field value map
    structItems, _ = reflections.Items(s)
```

##### Tags

*Tags* returns the structure's fields tag with the provided key. You can whether provide *Tags* a structure or a pointer to structure as first argument.


```go
    s := MyStruct {
        FirstField: "first value",      `matched:"first tag"`
        SecondField: 2,                 `matched:"second tag"`
        ThirdField: "third value",      `unmatched:"third tag"`
    }

    var structTags map[string]string

    // Tags will return a field name to tag content
    // map. Nota that only field with the tag name
    // you've provided which will be matched.
    // Here structTags will contain:
    // {
    //     "FirstField": "first tag",
    //     "SecondField": "second tag",
    // }
    structTags, _ = reflections.Tags(s, "matched")
```

#### Set a structure field value

*SetField* update's a structure's field value with the one provided. Note that
unexported fields cannot be set, and that field type and value type have to match.

```go
    s := MyStruct {
        FirstField: "first value",
        SecondField: 2,
        ThirdField: "third value",
    }

    // In order to be able to set the structure's values,
    // a pointer to it has to be passed to it.
    _ := reflections.SetField(&s, "FirstField", "new value")

    // If you try to set a field's value using the wrong type,
    // an error will be returned
    err := reflection.SetField(&s, "FirstField", 123)  // err != nil
```

## Important notes

* **unexported fields** cannot be accessed or set using reflections library: the golang reflect library intentionaly prohibits unexported fields values access or modifications.


## Contribute

* Check for open issues or open a fresh issue to start a discussion around a feature idea or a bug.
* Fork `the repository`_ on GitHub to start making your changes to the **master** branch (or branch off of it).
* Write tests which shows that the bug was fixed or that the feature works as expected.
* Send a pull request and bug the maintainer until it gets merged and published. :) Make sure to add yourself to AUTHORS_.

[the repository](http://github.com/oleiade/reflections)
[AUTHORS](https://github.com/oleiade/reflections/blob/master/AUTHORS.md)


[![Bitdeli Badge](https://d2weczhvl823v0.cloudfront.net/oleiade/reflections/trend.png)](https://bitdeli.com/free "Bitdeli Badge")

