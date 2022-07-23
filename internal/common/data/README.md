# Pagination

`/?$page=<int>&$size=<int>`

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
### Sort

`/?$sort=<field> [<order>],<field> [<order>]`
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

# RHS Colon
`/?date=$lt:10&date=$gte:2` will be: 
```
{
    date: {
        $lt: 10,
        $gte: 2
    }
}
```

## Filter (not implemented)
```
/?filter=<field> <op> <value>
/?filter=<field> <op> <value> and <field> <op> <value>
/?filter=<field> <op> <value> or <field> <op> <value>
/?filter=<field> <op> <value> or <field> <op> <value> and <field> <op> <value>
/?filter=<field> <op> <value> or (<field> <op> <value> and <field> <op> <value>)
```

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
        $in: ['ASSET', 'cost']
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
        $in: ['ASSET', 'cost']
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
                $in: ['ASSET', 'cost']
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
