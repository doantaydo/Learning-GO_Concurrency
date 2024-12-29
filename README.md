# GO Learning - Concurrency

<strong>Status:</strong> Learning ... <br />

This is the  repository for my learning course in Udemy.
- <strong>Course:</strong> [Working with Concurrency in Go (Golang)](https://gameloft.udemy.com/course/working-with-concurrency-in-go-golang)
- <strong>Instructor:</strong> [Ph.D. Trevor Sawler](https://www.udemy.com/course/building-modern-web-applications-with-go/#instructor-1)

In this course, I will:
- Look at the basic type in the sync package: mutexes (semaphores), and wait groups.
- We'll go through 3 of the classic computer science problems:
    + The Producer/Consumer problem
    + The Dining Philosopher problem
    + The Sleeping Barber problem
- Cover a more real-world scenario: Build a subset of a larger (imaginary) problem where a user of a service wants to subscribe and buy one of a series of available subscriptions. When user purchases a subscription, I will concurrently do:
    + Generate invoice
    + Send an email
    + Generate a PDF manual and send that to the user