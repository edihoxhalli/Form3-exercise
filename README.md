# Form3-exercise

### **Name:** Edi Hoxhalli
### **Exp:** No experience in developing in go for enterprise projects. Have been studying it intensively since Jan 2022.
<br>
<p>I think that this is just a fairly simple and concise implementation of an API client lib which hides away the implementation details and only exposes the three funcs for the required operations (create, read, delete). Also, certain variables are exported as well (Host, Version etc.), so that they can be changed from the caller as required, to keep up with the evolution of the api, and the multiple environments.</p>
<br>
<p>It has a unit test coverage of 96% and three simple integration tests, which check for expected responses' equality. More integration tests were considered for creation, but it was felt that they would be checking for validity of the actual (mock) Account API rather than the client library. </p>
<br>
<p>Certainly, there could be better alternatives for various parts of code encountered in the library. However, for the time being, this was considered good enough for commercial-ready status.</p>

## Thanks,
### Edi