# findemailaddrs

version 0.0.4 : 05 November 2023

Find email files on disk (by default `.eml` files) and parse the email
of each file and extract the sender, to and cc email addresses and save
these to a tab separated value (tsv) file removing duplicates.

## Example

```
./findemailaddrs -d <directory> -o <output> [-v] [-s ".ext"]

Look for files with the default ".eml" suffix or that optionally provided
with the "-s" flag in the directory rooted at <directory> and extract
the email addresses and associated names (where available) to <output>
in tab separated format.

Provide the -v flag for verbose output.

Options:
  -d string
    	path to directory to start eml file search
  -o string
    	file to save output
  -s string
    	file suffix to search for (default ".eml")
  -v	verbose
```

Example output with the -v flag:

```
./findemailaddrs -v -d /vols/archive -o report.tsv
...
processing directory: Juniper Berry Hotel
   file: Booking cancellation.eml
   file: Cancellation .eml
   file: Special Request for Reservation.eml
processing directory: BookShop
   file: BookShop Order No. XNY44897 Confirmed.eml
   file: BookShop Order No. XNY44897 Order Received.eml
counter 2962
unique addresses 719
```

## License

This project is licensed under the [MIT Licence](LICENCE).
