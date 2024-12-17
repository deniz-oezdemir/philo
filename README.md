# Philosophers in Go

## Learnings
- The `time.Sleep` function in Go is not highly precise, as its accuracy can be influenced by the operating system's scheduler and the system clock's granularity. This means it cannot guarantee nanosecond-level precision.

- My custom `preciseSleep` function offers improved accuracy but still depends on `time.Sleep`, making it less precise than its implementation in C (TODO: insert link).
