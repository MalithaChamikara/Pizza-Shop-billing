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

const Beverage = () => {
  const [beverages, setBeverages] = useState([]);
  const [error, setError] = useState("");
  const [open, setOpen] = useState(false);
  const [isEditing, setIsEditing] = useState(false);
  const [newBeverage, setNewBeverage] = useState({
    beverage_id: "",
    name: "",
    price: 0,
  });

  // Function to fetch beverage types from the backend using Axios
  const fetchBeverages = async () => {
    try {
      const response = await axios.get("http://localhost:8080/beverages"); // Replace with your backend endpoint
      if (Array.isArray(response.data)) {
        setBeverages(response.data);
      } else {
        setBeverages([]); // Ensure beverages is always an array
      }
    } catch (err) {
      console.error("Error fetching beverages:", err);
      setError("Failed to load beverages. Please try again later.");
    }
  };

  // Trigger API call on component mount
  useEffect(() => {
    fetchBeverages();
  }, []);
  // Function to handle opening the modal for adding a new beverage
  const handleOpen = () => {
    setIsEditing(false); // Set to false for creating a new beverage
    setNewBeverage({ beverage_id: "", name: "", price: 0 }); // Reset the form data
    setOpen(true);
  };
  // Function to handle opening the modal for editing a beverage
  const handleEditBeverage = (beverageId) => {
    const beverageToEdit = beverages.find(
      (bev) => bev.beverage_id === beverageId
    );
    if (beverageToEdit) {
      setIsEditing(true); // Set to true for editing
      setNewBeverage({ ...beverageToEdit }); // Set form fields with the existing beverage data
      setOpen(true);
    }
  };

  // Function to handle closing the modal
  const handleClose = () => {
    setOpen(false);
    setNewBeverage({ beverage_id: "", name: "", price: 0 }); // Reset form data
  };
  // Function to handle adding or updating a beverage
  const handleAddOrUpdateBeverage = async () => {
    try {
      const payload = {
        beverage_id: newBeverage.beverage_id.trim(), // Trim any extra spaces from id
        name: newBeverage.name,
        price: parseFloat(newBeverage.price), // Convert price to a decimal
      };

      // Validate input fields
      if (!payload.beverage_id || !payload.name || isNaN(payload.price)) {
        alert("Please provide a valid ID, name, or price.");
        return;
      }

      let response;
      if (isEditing) {
        // If editing, send PUT request to update the beverage
        response = await axios.put(
          `http://localhost:8080/beverages/${newBeverage.beverage_id}`,
          payload
        );
      } else {
        // If adding, send POST request to add the beverage
        response = await axios.post("http://localhost:8080/beverages", payload);
      }

      if (
        response.status === 201 ||
        response.status === 204 ||
        response.status === 200
      ) {
        alert(
          isEditing
            ? "Beverage updated successfully!"
            : "Beverage added successfully!"
        );
        fetchBeverages();
        handleClose();
      } else {
        alert("Failed to save beverage. Please try again.");
      }
    } catch (err) {
      console.error("Error saving beverage:", err);
      alert("Error saving beverage. Please try again.");
    }
  };

  // Function to handle deleting a beverage
  const handleDeleteBeverage = async (beverageId) => {
    if (
      window.confirm(
        `Are you sure you want to delete beverage with ID: ${beverageId}?`
      )
    ) {
      try {
        // Sending DELETE request to the backend to delete the beverage
        const response = await axios.delete(
          `http://localhost:8080/beverages/${beverageId}`
        );

        // Check if the response is successful (status 200-299)
        if (response.status >= 200 && response.status < 300) {
          // Alert user about successful deletion
          alert(`Deleted beverage with ID: ${beverageId}`);
          fetchBeverages();
        } else {
          // If the response is not successful, show an error
          alert("Failed to delete the beverage. Please try again.");
        }
      } catch (err) {
        // Catch any error that occurs during the API call
        console.error("Error deleting beverage:", err);
        alert("Error deleting the beverage. Please try again.");
      }
    }
  };

  return (
    <div>
      <h2>Beverages</h2>
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
        Add New Beverage
      </Button>
      {error ? (
        <p style={{ color: "red" }}>{error}</p>
      ) : beverages.length === 0 ? (
        <p>No beverages available.</p> // Message when no beverages are found
      ) : (
        <TableContainer component={Paper}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Beverage ID</TableCell>
                <TableCell>Beverage Name</TableCell>
                <TableCell>Price</TableCell>
                <TableCell>Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {beverages.map((beverage) => (
                <TableRow key={beverage.beverage_id}>
                  <TableCell>{beverage.beverage_id}</TableCell>
                  <TableCell>{beverage.name}</TableCell>
                  <TableCell>${beverage.price.toFixed(2)}</TableCell>
                  <TableCell>
                    {/* Edit and Delete Icons */}
                    <IconButton
                      onClick={() => handleEditBeverage(beverage.beverage_id)}
                      color="primary"
                    >
                      <EditIcon />
                    </IconButton>
                    <IconButton
                      onClick={() => handleDeleteBeverage(beverage.beverage_id)}
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
      {/* Modal for Adding or Editing Beverage */}
      <Dialog open={open} onClose={handleClose}>
        <DialogTitle>
          {isEditing ? "Edit Beverage" : "Add New Beverage"}
        </DialogTitle>
        <DialogContent>
          <DialogContentText>
            Please enter the required fields
          </DialogContentText>
          <TextField
            autoFocus
            margin="dense"
            label="Beverage ID"
            fullWidth
            required
            value={newBeverage.beverage_id}
            onChange={(e) =>
              setNewBeverage({ ...newBeverage, beverage_id: e.target.value })
            }
            disabled={isEditing} // Disable editing the ID when editing
          />
          <TextField
            autoFocus
            margin="dense"
            label="Beverage Name"
            fullWidth
            required
            value={newBeverage.name}
            onChange={(e) =>
              setNewBeverage({ ...newBeverage, name: e.target.value })
            }
          />
          <TextField
            margin="dense"
            label="Price"
            type="number"
            fullWidth
            required
            value={newBeverage.price}
            onChange={(e) =>
              setNewBeverage({ ...newBeverage, price: e.target.value })
            }
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClose} color="secondary">
            Cancel
          </Button>
          <Button onClick={handleAddOrUpdateBeverage} color="primary">
            {isEditing ? "Update" : "Add"}
          </Button>
        </DialogActions>
      </Dialog>
    </div>
  );
};

export default Beverage;
