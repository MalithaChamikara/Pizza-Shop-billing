import React, { useState, useEffect } from "react";
import axios from "axios";
import {
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Button,
  IconButton,
  Modal,
  Box,
  TextField,
  MenuItem,
  Select,
  InputLabel,
  FormControl,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
} from "@mui/material";
import EditIcon from "@mui/icons-material/Edit";
import DeleteIcon from "@mui/icons-material/Delete";
const Pizzatype = () => {
  const [pizzaTypes, setPizzaTypes] = useState([]);
  const [toppings, setToppings] = useState({});
  const [selectedToppings, setSelectedToppings] = useState([]);
  const [open, setOpen] = useState(false);
  const [newPizza, setNewPizza] = useState({
    pizza_type_id: "",
    name: "",
    size: "",
    base_price: 0,
    description: "",
  });
  const [error, setError] = useState("");

  // Define available sizes
  const pizzaSizes = ["Small", "Medium", "Large"];

  // Function to fetch pizza types from the backend using Axios
  const fetchPizzaTypes = async () => {
    try {
      const response = await axios.get("http://localhost:8080/pizzas");
      console.log(response.data);
      if (Array.isArray(response.data)) {
        setPizzaTypes(response.data);
      } else {
        setPizzaTypes([]);
      }
    } catch (err) {
      console.error("Error fetching pizza types:", err);
      setError("Failed to load pizza types. Please try again later.");
    }
  };
  // Function to fetch toppings for a specific pizza type
  const fetchToppingsForPizzaType = async (pizzaTypeId) => {
    try {
      const response = await axios.get(
        `http://localhost:8080/pizzas/${pizzaTypeId}/toppings`
      );
      console.log(response.data);
      return response.data; // Return topping names as an array
    } catch (err) {
      console.error(
        `Error fetching toppings for pizza type ${pizzaTypeId}:`,
        err
      );
      return [];
    }
  };

  // Fetch toppings for all pizza types
  const fetchAllToppings = async () => {
    const toppingMap = {};
    for (const pizza of pizzaTypes) {
      const toppingNames = await fetchToppingsForPizzaType(pizza.pizza_type_id);
      toppingMap[pizza.pizza_type_id] = toppingNames;
    }
    setToppings(toppingMap); // Store toppings in state
  };

  // Fetch data on component mount
  useEffect(() => {
    const fetchData = async () => {
      await fetchPizzaTypes();
    };
    fetchData();
  }, []);

  // Fetch toppings when pizza types are loaded
  useEffect(() => {
    if (pizzaTypes.length > 0) {
      fetchAllToppings();
    }
  }, [pizzaTypes]);

  // Handle opening modal for adding new pizza
  const handleOpenModal = () => {
    setOpen(true);
  };

  // Handle closing modal
  const handleCloseModal = () => {
    setOpen(false);
  };

  // Handle change in pizza input fields
  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setNewPizza({ ...newPizza, [name]: value });
  };
  // Handle change in size select box
  const handleSizeChange = (event) => {
    setNewPizza({ ...newPizza, size: event.target.value });
  };

  // Function to handle adding a new pizza
  const handleAddPizza = async () => {
    console.log(newPizza);
    try {
      const response = await axios.post(
        "http://localhost:8080/pizzas",
        newPizza
      );
      if (response.status === 201) {
        alert("Pizza added successfully!");
        fetchPizzaTypes();
        setOpen(false); // Close the modal after adding pizza
      } else {
        alert("Failed to add pizza. Please try again.");
      }
    } catch (err) {
      console.error("Error adding pizza:", err);
      alert("Error adding pizza. Please try again.");
    }
  };

  // Function to handle editing a pizza
  const handleEditPizza = (pizzaId) => {
    alert(`Edit pizza with ID: ${pizzaId}`);
  };

  // Function to handle deleting a pizza
  const handleDeletePizza = async (pizzaId) => {
    if (
      window.confirm(
        `Are you sure you want to delete pizza with ID: ${pizzaId}?`
      )
    ) {
      try {
        // Sending DELETE request to the backend to delete the pizzatype

        const response = await axios.delete(
          `http://localhost:8080/pizzas/${pizzaId}`
        );

        // Check if the response is successful (status 200-299)
        if (response.status >= 200 && response.status < 300) {
          // Alert user about successful deletion
          alert(`Deleted pizza with ID: ${pizzaId}`);
          fetchPizzaTypes();
        } else {
          // If the response is not successful, show an error
          alert("Failed to delete the pizza. Please try again.");
        }
      } catch (err) {
        // Catch any error that occurs during the API call
        console.error("Error deleting pizza:", err);
        alert("Error deleting the pizza. Please try again.");
      }
    }
  };

  return (
    <div>
      <h2>Pizza Types</h2>
      <Button
        variant="contained"
        color="primary"
        onClick={handleOpenModal}
        style={{
          position: "absolute",
          top: 20,
          right: 20,
        }}
      >
        Add New Pizza
      </Button>
      {error ? (
        <p style={{ color: "red" }}>{error}</p>
      ) : pizzaTypes.length === 0 ? (
        <p>No pizzas available.</p> // Message when no beverages are found
      ) : (
        <TableContainer component={Paper}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Pizza ID</TableCell>
                <TableCell>Pizza Name</TableCell>
                <TableCell>Size</TableCell>
                <TableCell>Price</TableCell>
                <TableCell>Description</TableCell>
                <TableCell>Toppings</TableCell>
                <TableCell>Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {pizzaTypes.map((pizza) => (
                <TableRow key={pizza.pizza_type_id}>
                  <TableCell>{pizza.pizza_type_id}</TableCell>
                  <TableCell>{pizza.name}</TableCell>
                  <TableCell>{pizza.size}</TableCell>
                  <TableCell>${pizza.base_price.toFixed(2)}</TableCell>
                  <TableCell>{pizza.description}</TableCell>
                  <TableCell>
                    {Array.isArray(toppings[pizza.pizza_type_id])
                      ? toppings[pizza.pizza_type_id].join(", ")
                      : "Loading..."}
                  </TableCell>
                  <TableCell>
                    {/* Edit and Delete Icons */}
                    <IconButton
                      onClick={() => handleEditPizza(pizza.pizza_type_id)}
                      color="primary"
                    >
                      <EditIcon />
                    </IconButton>
                    <IconButton
                      onClick={() => handleDeletePizza(pizza.pizza_type_id)}
                      color="secondary"
                    >
                      <DeleteIcon />
                    </IconButton>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      )}
      {/* Modal for adding new pizza */}
      <Dialog open={open} onClose={handleCloseModal}>
        <DialogTitle>Add New Pizza</DialogTitle>
        <DialogContent>
          <TextField
            label="Pizza ID"
            name="pizza_type_id"
            value={newPizza.pizza_type_id}
            onChange={handleInputChange}
            fullWidth
            margin="normal"
          />
          <TextField
            label="Pizza Name"
            name="name"
            value={newPizza.name}
            onChange={handleInputChange}
            fullWidth
            margin="normal"
          />

          <FormControl fullWidth margin="normal">
            <InputLabel>Size</InputLabel>
            <Select
              value={newPizza.size}
              onChange={handleSizeChange}
              label="Size"
              name="size"
            >
              {pizzaSizes.map((size) => (
                <MenuItem key={size} value={size}>
                  {size}
                </MenuItem>
              ))}
            </Select>
          </FormControl>

          <TextField
            label="Base Price"
            name="base_price"
            type="number"
            value={newPizza.base_price}
            onChange={handleInputChange}
            fullWidth
            margin="normal"
          />
          <TextField
            label="Description"
            name="description"
            value={newPizza.description}
            onChange={handleInputChange}
            fullWidth
            margin="normal"
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseModal} color="primary">
            Cancel
          </Button>
          <Button onClick={handleAddPizza} color="primary">
            Add Pizza
          </Button>
        </DialogActions>
      </Dialog>
    </div>
  );
};

export default Pizzatype;
