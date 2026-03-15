# Title 1

## Table of contents

- [Title 2](#title-2)
  - [Title 3. With `a code` word](#title-3-with-a-code-word)
- [Text](#text)

## Title 2

### Title 3. With `a code` word

## Text

Foo bar.

Multi line code:

```
print('foo')

# Comment here.
print('bar')
```

Multi line code with code type:

```bash
foo
```

This is `a piece of code` and this is `another one`.

Link to a [website](https://cmoli.es).

Another web: <https://cmoli.es>.

[This line only contains a link](https://cmoli.es)

Link to a [file](../testdata/file.html).

[This line only contains a link to a file](../testdata/file.html)

The following line is another one with only a link:

<https://cmoli.es>

List:

1. a
    - a.1

      a.1.1

    - a.2

      a.2.2

    a.3

1. b

    Code in list:

    ```
    foo
    ```

Table:

column 1 | column 2
---------|---------
foo foo  | bar bar
baz baz  | qux qux

This line has the character used in tables: |. But is not a table.

A blockquote:

> This is a blockquote with an [url](https://cmoli.es) and `a piece of code`.

A blockquote with indentation:

- List element:

  > I am an indented blockquote.

This is a line with the blockquote character > but this is not a blockquote.

Code and blockquote:

```
foo # <blockquote character and close character too>
bar # <blockquote character
baz # line after other line with open blockquote character
```

Inline code with blockquote character: `foo <bar`.

Image:

![](favicon.ico)

Image with alt:

![with alt](favicon.ico)

## Underscore: \_\_init__ \_\_init\_\_

Italics: _init_.

Bold: __init__.

Escaped underscore: \_\_init\_\_.

Enough to escape first underscores: \_\_init__.

Multi line code with underscore:

```
__init__
```

Inline code with underscore: `__init__`.

- List element with escaped underscore: \_\_init\_\_.
- Escape only first underscores: \_\_init__.
