"use client";
import React from "react";
import { Form, Input, Button } from "@heroui/react";

export const animals = [
  { key: "1", label: "Cat" },
  { key: "2", label: "Dog" },
  { key: "4", label: "Elephant" },
  { key: "3", label: "Lion" },
  { key: "tiger", label: "Tiger" },
  { key: "giraffe", label: "Giraffe" },
  { key: "dolphin", label: "Dolphin" },
  { key: "penguin", label: "Penguin" },
  { key: "zebra", label: "Zebra" },
  { key: "shark", label: "Shark" },
  { key: "whale", label: "Whale" },
  { key: "otter", label: "Otter" },
  { key: "crocodile", label: "Crocodile" },
];

function LoginForm() {
  return (
    <Form className="w-full max-w-xs flex flex-col gap-3">
      <Input
        label="Username"
        labelPlacement="outside"
        name="username"
        placeholder="Enter your username"
      />
      <Button type="submit" variant="flat">
        Submit
      </Button>
    </Form>
  );
}

export default LoginForm;
