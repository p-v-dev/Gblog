# Future improvements

- [ ] Add `User user.User` BelongsTo field to BlogPost (allows `post.User.Name` without extra query)
- [ ] Add `/users/:id/posts` route to list posts by user
- [ ] Validate UserID exists before creating post (app-level check, not just DB FK)
- [ ] Add `/api/v1/users` CRUD routes (GET, PUT, DELETE, list)
