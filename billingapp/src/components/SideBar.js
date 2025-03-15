import React from "react";
import { List, ListItem, ListItemText, Paper } from "@mui/material";
import { Link } from "react-router-dom";

const SideBar = () => {
  return (
    <Paper
      elevation={3} // Shadow effect
      sx={{
        height: "100vh", // Full viewport height
        width: "100%", // Take full width of the Grid container
        padding: 2,
      }}
    >
      <List>
        <ListItem button component={Link} to="/pizzatypes">
          <ListItemText primary="Pizza Types" />
        </ListItem>
        <ListItem button component={Link} to="/toppings">
          <ListItemText primary="Toppings" />
        </ListItem>
        <ListItem button component={Link} to="/beverages">
          <ListItemText primary="Beverages" />
        </ListItem>
        <ListItem button component={Link} to="/invoices">
          <ListItemText primary="Invoice" />
        </ListItem>
      </List>
    </Paper>
  );
};

export default SideBar;
