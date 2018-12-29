# Testing zoo

This is a directory that contains a set of test cases that exhibit edge cases
and enshrine expected functionality in code. These tests can be invoked by
`go test`ing the `cmd` directory.


## How do they work?

These tests simulate running the logbook generator using the name "Andrew
Allen", since it is the best name, and invoke the log entry generator once for
each of the files that end in `.out` using the remainder of the file name as
the date to generate. For example, let's look at the `simple` test case that
exists.

```
└── simple
    ├── 2000-01-01.md
    └── 2000-01-02.out
```

The contents of `2000-01-01.md` are:

```

# Andrew Allen - 2000-01-01

Sample first entry

Tomorrow: Do stuff.
```

The contents of `2000-01-02.out` are:

```
# Andrew Allen - 2000-01-02

Reminders:
From 2000-01-01: Tomorrow text
```

The test case performs the following actions:

 *  Reads in the list of files
 *  Determines that you're interested in generating an entry for `2000-01-02`
 *  Copies all the files into a temporary directory to be non-destructive.
 *  Parses all previous entries (only `2000-01-01.md` in this case)
 *  Generates an `2000-01-02.md` and compares it to `2000-01-02.out` and if
    they are equivilent, the test passes.

## Writing/updating a test cases

Writing and updating a test case to demonstrate a bug, or exercise a new
feature should be fairly simple and should, most of the time, not require
writing any code.

 *  Make a new directory in the `testdata` directory with a descriptive name
    of what you would like to demonstrate.
 *  Then fill the file with some sample `.md` files, named in accordance with
    the scheme.
 *  Create an empty `.out` file for each of the dates you would like to
    generate.
 *  Now invoke the test with the `--update_goldens` flag.

    ```
    go test github.com/achew22/logbook/cmd --update_goldens
    ```

    This will fill in the `.out` files with what is generated today.
 *  Run `git commit` to lock your work in place.
 *  If you are trying to demonstrate a bug, update the `.out` file to reflect
    what you expect the output should have been.
 *  Run `git commit` to create a simple diff of what it does now vs what your
    expectations are.
 *  Upload to GitHub and mail the PR.
