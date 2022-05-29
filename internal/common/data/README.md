# Pagination

`/?page=<int>&size=<int>`

## Request
### interface `Pageable`
``` go
type Pageable interface {
	getPage() int
	getSize() int
	getOffset() int64
	getSort() []Sort
}
```
### Sort (not implemented)

`/?sort=<field> [<order>],<field> [<order>]`
```
[
    { <field>: DESC|ASC },
    { <field>: DESC|ASC }
]
```

## Response
### struct `Page`
```
{
    content: [<entries>],
    page: <int>,
    total: <int>,
    size: <int>,
    count: <int>,
    isFirst: <bool>,
    isLast: <bool>
}
```

# Selector
## Choose

`/?choose=<field1>,<field2>`

## Find (not implemented)

`/?find=<field> <op> <value> and|or <field> <op> <value>`

**equal:**
```
{
    <field>: "<value>"
}
```
e.g.:
```
{
    title: "something"
}
```

**operator query:**
```
{
    <field>: {
        <operator>: <value>|<value>[]
    }
}
```
e.g.:
```
{
    accountType: {
        $in: ['ASSET', 'COST']
    }
}
```

**and condition:**
```
{
    <field>: <operators>|<value>,
    <field>: <operators>|<value>,
    ...
}
```
e.g.:
```
{
    title: "something",
    accountType: {
        $in: ['ASSET', 'COST']
    }
}
```

**or condition:**
```
{
    $or: [
        {<find>},
        {<find>}
    ]
}
```
e.g.:
```
{
    $or: [
        { title: "something" },
        {
            accountType: {
                $in: ['ASSET', 'COST']
            }
        }
    ]
}
```

**nested element:**
```
// TODO
```

**array:**
```
// TODO
```
