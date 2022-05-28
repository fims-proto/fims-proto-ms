# Search:

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

# Page

`/?page=<int>&size=<int>`