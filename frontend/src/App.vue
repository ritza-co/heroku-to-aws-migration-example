<script >
import AddEmployeeForm from './components/AddEmployeeForm.vue';
export default {
    data() {
        return {
            employees: [],
        };
    },
    mounted() {
        this.fetchEmployeeData();
    },
    methods: {
        fetchEmployeeData() {
            const apiUrl = '/employees';
            // Using Axios to fetch data
            this.axios
                .get(apiUrl)
                .then((response) => {
                this.employees = response.data;
            })
                .catch((error) => {
                console.error('Error fetching employee data:', error);
            });
        },
        deleteEmployee(id) {
            const apiUrl = import.meta.env.VITE_API_URL + '/employees/' + id;
            // Using Axios to fetch data
            this.axios
                .delete(apiUrl)
                .then((response) => {
                this.employees = response.data;
                location.reload();
            })
                .catch((error) => {
                console.error('Error fetching employee data:', error);
            });
        },
    },
    components: { AddEmployeeForm }
};
</script>

<template>
  <header>
    <h1>Employee Data - Acme Inc</h1>
  </header>
  <main>
    <AddEmployeeForm />
    <h2>Current Employees</h2>
    <table>
      <thead>
        <tr>
          <th>ID</th>
          <th>First Name</th>
          <th>Last Name</th>
          <th>Email</th>
          <th>Hourly Pay</th>
          <th>Current Earnings</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="employee in employees" :key="employee._id">
          <td>{{ employee._id }}</td>
          <td>{{ employee.firstname }}</td>
          <td>{{ employee.lastname }}</td>
          <td>{{ employee.email }}</td>
          <td>{{ employee.hourlypay }}</td>
          <td>{{ employee.currentearnings }}</td>
          <th><button @click="() => deleteEmployee(employee._id)">Delete</button></th>

        </tr>
      </tbody>
    </table>
  </main>
</template>

<style scoped>
header {
  line-height: 1.5;
}

.logo {
  display: block;
  margin: 0 auto 2rem;
}

@media (min-width: 1024px) {
  header {
    display: flex;
    place-items: center;
    padding-right: calc(var(--section-gap) / 2);
  }

  .logo {
    margin: 0 2rem 0 0;
  }

  header .wrapper {
    display: flex;
    place-items: flex-start;
    flex-wrap: wrap;
  }
}
table {
  width: 100%; /* Makes the table expand to the full width of its container */
  border-collapse: collapse; /* Removes the space between table cells */
}

th, td {
  text-align: left; /* Aligns text to the left */
  padding: 8px; /* Adds padding inside each cell */
  border: 1px solid #ddd; /* Adds a border to each cell */
}

th {
  background-color: #f4f4f4; /* Adds a background color to the header cells */
}

tr:nth-child(even) {
  background-color: #f9f9f9; /* Adds a striped effect to the table rows */
}
</style>
