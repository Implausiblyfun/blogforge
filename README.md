# blogforge

Blogforge is a simple timeboxed attempt at a basic blog backend.

It is not intended for scale and is essentially a set of endpoints that power a basic CRUD setup for blog posts scoped to users and a synonomous setup for the creation of the aformentioned users.


Long story short I was doing some coding application thing and started pulling in random silly stuff.

On the note of reuse, I havent spent the time to collate the most permissive total license on this but I have done what I think is a relatively decent job of denoting when I shamelessly steal an example.
Given this if I have time I will find the maximally permissive acceptable license but in the meantime use snippets at your own risk as all snippets are acceptable for non-coomericial use 
but are not gauranteed for commercial uses (not that any of this should be used in production as it most definitely SHOULD NOT BE)

Once again this is not a good example of go or blogs or anything really but I figured it existing on github wasnt the worst coding ive ever done.

So here it goes.



## Running the thing
Have docker
Dont have 8080 already mapped or 3306 I guess
Then run:
```
docker-compose -f "docker-compose.yml" up -d --build
```

then do somethings like
```
curl localhost:8080/index
```

// Now attempt the create user route with the second attempt working to create usr1
```
curl -i -X POST -H "Content-Type: application/json" -d "{\"username\":\"notgoingtowork\"}" http://localhost:8080/newuser
curl -i -X POST -H "Content-Type: application/json" -d "{\"username\": \"username\", \"password\": \"password\", \"first_name\": \"nate\", \"last_name\": \"food\"}" http://localhost:8080/newuser
curl -i -X POST -H "Content-Type: application/json" -d "{\"username\": \"usr1\", \"password\": \"12345\", \"first_name\": \"nate\", \"last_name\": \"food\"}" http://localhost:8080/newuser
```

and then nothing else real is implemented
so i guess open up the jeager url to see some of the gin interactions as they are instrumented lightly so thats cool