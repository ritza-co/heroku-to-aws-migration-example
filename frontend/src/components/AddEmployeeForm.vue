<script>
export default {
    name: 'PostFormAxios',
    data(){
        return{
            form: {
                firstname: '',
                lastname: '',
                email: '',
                hourlypay: '',
            }
        }
    },
    methods:{
        submitForm(){
          const apiUrl = '/employees';
            // Using Axios to fetch data
          const url = "https://cors-anywhere.herokuapp.com/" + import.meta.env.VITE_API_URL + apiUrl

            this.axios.post(url, this.form)
                 .then((res) => {
                     //Perform Success Action
                     alert("Employee Added Successfully")
                 })
                 .catch((error) => {
                     // error.response.status Check status code
                     console.log(error.response.data)
                 }).finally(() => {
                     //Perform action in always
                     location.reload()
                 });
        }
    }
}
</script>
<template>
  <div>
          <h2> Add an employee </h2>
          <form v-on:submit.prevent="submitForm">
              <div class="form-group">
                  <label for="firstName">First Name</label>
                  <input type="text" class="form-control" id="firstName" placeholder="First name" v-model="form.firstname">
              </div>
              <div class="form-group">
                  <label for="lastName">Last Name</label>
                  <input type="text" class="form-control" id="lastName" placeholder="Last name" v-model="form.lastname">
              </div>
              <div class="form-group">
                  <label for="email">Email Address</label>
                  <input type="email" class="form-control" id="email" placeholder="name@example.com"
                      v-model="form.email">
              </div>
              <div class="form-group">
                  <label for="hourlyPay">Hourly Pay</label>
                  <input type="hourlyPay" class="form-control" id="hourlyPay" placeholder="250"
                      v-model="form.hourlypay">
              </div>
              <div class="form-group">
                  <button class="btn btn-primary">Submit</button>
              </div>
          </form>
  </div>
  <hr />
</template>

<style scoped>
.form-group {
  margin-bottom: 15px; /* Adds space between form groups */
}

label {
  display: block; /* Makes the label take the full width */
  margin-bottom: 5px; /* Adds a small space between the label and the input field */
  font-weight: bold; /* Makes the label text bold */
}

input[type="text"],
input[type="email"],
input[type="hourlyPay"] {
  width: 100%; /* Makes the input fields take the full width */
  padding: 8px; /* Adds padding inside the input fields for better readability */
  border: 1px solid #ddd; /* Adds a subtle border to the input fields */
  border-radius: 4px; /* Rounds the corners of the input fields */
  box-sizing: border-box; /* Ensures padding does not affect the overall width */
}

input[type="text"]:focus,
input[type="email"]:focus,
input[type="hourlyPay"]:focus {
  border-color: #0056b3; /* Changes the border color when the input field is focused */
  outline: none; /* Removes the default focus outline */
}

.btn-primary {
  background-color: #007bff; /* Primary button background color */
  color: white; /* Text color for the primary button */
  padding: 10px 15px; /* Padding inside the button */
  border: none; /* Removes the border */
  border-radius: 4px; /* Rounds the corners of the button */
  cursor: pointer; /* Changes the cursor to a pointer when hovering over the button */
}

.btn-primary:hover {
  background-color: #0056b3; /* Darkens the button when hovered */
}
</style>