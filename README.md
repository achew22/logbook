# Logbook

There are people who are good at keeping logs of what they have done through the
day. I'm not really one of them, but I figured if I had a nice tool that helped
me remember some stuff it would make it easier.

This program creates a file in the "logbook" folder of my home directory with
today's date. It then scans old entries from the past to see if I left any
remarks to be reminded of on today's date and puts them in the list. Then it
exits.

## Example output

If you run this program on 2018-01-01, it will make a file
`$HOME/logbook/2018-01-01.md` with the following content:

```
# Andrew Allen - 2018-01-01


```

As your day progressed you might make the file look like this;

```
# Andrew Allen - 2018-01-01

Today I did stuff and things.
Tomorrow: I will finish the thing I forgot to do.
```

Then when you ran the program on 2018-01-02, you would get the output:

```
# Andrew Allen - 2018-01-01

Reminders:
2018-01-01: I will finish the thing I forgot to do.


```

This provides a simple way to leave notes for yourself going forward in a place
you already use.

THIS IS NOT AN OFFICIAL GOOGLE PRODUCT.
