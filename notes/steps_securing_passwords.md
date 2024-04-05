
# Steps for Securing Passwords

Now that we have spent some time learning about SQL we are eager and ready to start adding users to our application, but unfortunately we aren’t quite ready to do that.

If we look at what information we need for a user, we will find that users are only composed of a few pieces of information:

- An ID
- An email address
- A password

Only three columns of data, simple, right?

Unfortunately, users are far from simple due to the password field. Unlike most data in our database, passwords cannot simply be stored as text in our database. If we were to do this, a data breach could result in every user’s password being leaked. Not only would this cause trouble with our own app, but many users share passwords between websites and this could lead to their accounts being compromised in other applications as well. As a result, our authentication setup is arguably the most important and sensitive part of our application.

While this might seem scary at first, the truth is that implementing a secure authentication service isn’t incredibly challenging, but it does require us to follow a series of industry standards and to *not deviate from those standards at all*. In many cases, security issues stem from developers with good intentions attempting to do something custom and inadvertently introducing a security issue. We don’t want to see our website on the front page of the New York Times with a “Data Breach” headline, so let’s stick to the industry standards.

Throughout the rest of this section we will go into many of these in detail, but for now let’s take a high level look at the steps to securely handling passwords:

1. Use HTTPs to secure our domain.
2. Store hashed passwords. Never store encrypted or plaintext passwords.
3. Add a salt to passwords before hashing.
4. Using time-constant functions during authentication.

The first one is pretty straightforward; when we deploy our application, we should get an SSL certificate and ensure that all traffic uses URLs with the `https` prefix. This ensures that all communications with our web server is encrypted, making it nearly impossible for someone to intercept the password when a user submits a form or does something similar.

We won’t actually look at setting up an SSL certificate until later in the course when we deploy to production, but it is worth noting here as a requirement for securely handling passwords.

The next requirement will be covered in great detail, but it can never stated enough. NEVER under any circumstances store an encrypted or plaintext password. Always store a hashed value instead. If you are ever unsure of whether something is a hash, encryption, or something similar, the simple rule is that if you can calculate the password from whatever you are storing, it is NOT a hash. We will go into this in detail in the next lesson as we learn what hash functions are.

Third, passwords should always use a salt. This is done to prevent [rainbow tables](https://www.beyondidentity.com/glossary/rainbow-table-attack) from being effective. We will explore this in detail in a future lesson.

The final one is a bit trickier to explain and understand, but hackers can actually use obscure information like how long it takes our server to check a password to try to determine which hashing function we are using. To prevent this, we will learn how to use time-constant functions when authenticating users.

