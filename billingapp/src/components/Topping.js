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
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
  TextField,
} from "@mui/material";
import EditIcon from "@mui/icons-material/Edit";
import DeleteIcon from "@mui/icons-material/Delete";

const Topping = () => {
  const [toppings, setToppings] = useState([]);
  const [error, setError] = useState("");
  const [open, setOpen] = useState(false);
  const [isEditing, setIsEditing] = useState(false);
  const [newTopping, setNewTopping] = useState({
    topping_id: "",
    name: "",
    price: 0,
  });

  // Function to fetch topping types from the backend using Axios
  const fetchToppings = async () => {
    try {
      const response = await axios.get("http://localhost:8080/toppings"); // Replace with your backend endpoint
      if (Array.isArray(response.data)) {
        setToppings(response.data);
      } else {
        setToppings([]); // Ensure toppings is always an array
      }
    } catch (err) {
      console.error("Error fetching toppings:", err);
      setError("Failed to load toppings. Please try again later.");
    }
  };

  // Trigger API call on component mount
  useEffect(() => {
    fetchToppings();
  }, []);

  // Function to handle opening the modal for adding a new topping
  const handleOpen = () => {
    setIsEditing(false); // Set to false for creating a new topping
    setNewTopping({ topping_id: "", name: "", price: 0 }); // Reset the form data
    setOpen(true);
  };

  // Function to handle opening the modal for editing a topping
  const handleEditTopping = (toppingId) => {
    const toppingToEdit = toppings.find((top) => top.topping_id === toppingId);
    if (toppingToEdit) {
      setIsEditing(true); // Set to true for editing
      setNewTopping({ ...toppingToEdit }); // Set form fields with the existing topping data
      setOpen(true);
    }
  };

  // Function to handle closing the modal
  const handleClose = () => {
    setOpen(false);
    setNewTopping({ topping_id: "", name: "", price: 0 }); // Reset form data
  };
  // Function to handle adding or updating a topping
  const handleAddOrUpdateTopping = async () => {
    try {
      const payload = {
        topping_id: newTopping.topping_id.trim(), // Trim any extra spaces from id
        name: newTopping.name,
        price: parseFloat(newTopping.price), // Convert price to a decimal
      };

      // Validate input fields
      if (!payload.topping_id || !payload.name || isNaN(payload.price)) {
        alert("Please provide a valid ID, name, or price.");
        return;
      }

      let response;
      if (isEditing) {
        // If editing, send PUT request to update the topping
        response = await axios.put(
          `http://localhost:8080/toppings/${newTopping.topping_id}`,
          payload
        );
      } else {
        // If adding, send POST request to add the topping
        response = await axios.post("http://localhost:8080/toppings", payload);
      }

      if (
        response.status === 201 ||
        response.status === 204 ||
        response.status === 200
      ) {
        alert(
          isEditing
            ? "Topping updated successfully!"
            : "Topping added successfully!"
        );
        fetchToppings();
        handleClose();
      } else {
        alert("Failed to save topping. Please try again.");
      }
    } catch (err) {
      console.error("Error saving topping:", err);
      alert("Error saving topping. Please try again.");
    }
  };

  // Function to handle deleting a topping
  const handleDeleteTopping = async (toppingId) => {
    if (
      window.confirm(
        `Are you sure you want to delete topping with ID: ${toppingId}?`
      )
    ) {
      try {
        // Sending DELETE request to the backend to delete the topping
        const response = await axios.delete(
          `http://localhost:8080/toppings/${toppingId}`
        );

        // Check if the response is successful (status 200-299)
        if (response.status >= 200 && response.status < 300) {
          // Alert user about successful deletion
          alert(`Deleted topping with ID: ${toppingId}`);
          fetchToppings();
        } else {
          // If the response is not successful, show an error
          alert("Failed to delete the topping. Please try again.");
        }
      } catch (err) {
        // Catch any error that occurs during the API call
        console.error("Error deleting topping:", err);
        alert("Error deleting the topping. Please try again.");
      }
    }
  };

  return (
    <div>
      <h2>Toppings</h2>
      <Button
        variant="contained"
        color="primary"
        onClick={handleOpen}
        style={{
          position: "absolute",
          top: 20,
          right: 20,
        }}
      >
        Add New Topping
      </Button>
      {error ? (
        <p style={{ color: "red" }}>{error}</p>
      ) : toppings.length === 0 ? (
        <p>No toppings available.</p> // Message when no toppings are found
      ) : (
        <TableContainer component={Paper}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Topping ID</TableCell>
                <TableCell>Topping Name</TableCell>
                <TableCell>Price</TableCell>
                <TableCell>Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {toppings.map((topping) => (
                <TableRow key={topping.topping_id}>
                  <TableCell>{topping.topping_id}</TableCell>
                  <TableCell>{topping.name}</TableCell>
                  <TableCell>${topping.price.toFixed(2)}</TableCell>
                  <TableCell>
                    {/* Edit and Delete Icons */}
                    <IconButton
                      onClick={() => handleEditTopping(topping.topping_id)}
                      color="primary"
                    >
                      <EditIcon />
                    </IconButton>
                    <IconButton
                      onClick={() => handleDeleteTopping(topping.topping_id)}
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
      {/* Modal for Adding or Editing Topping */}
      <Dialog open={open} onClose={handleClose}>
        <DialogTitle>
          {isEditing ? "Edit Topping" : "Add New Topping"}
        </DialogTitle>
        <DialogContent>
          <DialogContentText>
            Please enter the required fields
          </DialogContentText>
          <TextField
            autoFocus
            margin="dense"
            label="Topping ID"
            fullWidth
            required
            value={newTopping.topping_id}
            onChange={(e) =>
              setNewTopping({ ...newTopping, topping_id: e.target.value })
            }
            disabled={isEditing} // Disable editing the ID when editing
          />
          <TextField
            autoFocus
            margin="dense"
            label="Topping Name"
            fullWidth
            required
            value={newTopping.name}
            onChange={(e) =>
              setNewTopping({ ...newTopping, name: e.target.value })
            }
          />
          <TextField
            margin="dense"
            label="Price"
            type="number"
            fullWidth
            required
            value={newTopping.price}
            onChange={(e) =>
              setNewTopping({ ...newTopping, price: e.target.value })
            }
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClose} color="secondary">
            Cancel
          </Button>
          <Button onClick={handleAddOrUpdateTopping} color="primary">
            {isEditing ? "Update" : "Add"}
          </Button>
        </DialogActions>
      </Dialog>
    </div>
  );
};

export default Topping;
