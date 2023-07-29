PINGO(1)                  General Commands Manual                 PINGO(1)

NAME
     pingo – Guess the operating system of a network host

SYNOPSIS
     pingo target [--retries N] [--tolerance N] [--verbose]

DESCRIPTION
     The pingo command sends an ICMP echo request to a network host and uses
     the TTL (Time To Live) value in the echo reply along with the count of
     network hops to guess the host's operating system.

     The target is the hostname or IP address of the network host to ping.

OPTIONS
     --retries N
             Number of tries to ping the target and guess the operating
             system. Default is 1.

     --tolerance N
             Tolerance value for guessing the operating system based on TTL.
             If the difference between the TTL value from the echo reply and
             the typical TTL value of an operating system is within the
             tolerance, the operating system is considered a match. Default
             is 0.

     --verbose
             Enable verbose output. When enabled, the outputs of the
             underlying traceroute and ping commands are printed.

EXIT STATUS
     The pingo utility exits 0 on success, and >0 if an error occurs.

EXAMPLES
     Guess the operating system of the host at example.com:

           pingo example.com

     Make 5 tries to guess the operating system of the host at example.com:

           pingo example.com --retries 5

     Guess the operating system of the host at example.com with a TTL
     tolerance of 2:

           pingo example.com --tolerance 2

     Guess the operating system of the host at example.com with verbose
     output:

           pingo example.com --verbose

SEE ALSO
     ping(1), traceroute(1)

STANDARDS
     The pingo command is not a standard UNIX command and may not be
     available on all systems.

BUGS
     The operating system guess is based on typical TTL values and may not
     always be accurate.

AUTHORS
     Developed by hollerith.

hollerith                          July 29, 2023                      PINGO(1)
