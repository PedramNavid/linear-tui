## API Best Practices

Do's:

Register a programmatic webhook and get updates for all issues for the team. When you detect changes, update the issue information. You can also automatically register webhooks for OAuth applications.
If you have to poll recent changes, order results by returning recently updated issue first. See Pagination section above how to implement this
Filter issues in your GraphQL request instead of fetching all issues and filtering in code.
Dont's:

Poll updates for each issue in the application. There should never be a reason to do this and your application might get rate limited. See above tactics to implement this better
