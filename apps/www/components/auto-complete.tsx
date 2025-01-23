import React, { Key } from "react";
import {
  Autocomplete as HeroAutocomplete,
  AutocompleteItem,
  Input,
} from "@heroui/react";

export type FieldState = {
  selectedKey: React.Key | null;
  inputValue: string;
  items: typeof animals;
};

export const animals = [
  {
    label: "Cat",
    key: "cat",
    description: "The second most popular pet in the world",
  },
  {
    label: "Dog",
    key: "dog",
    description: "The most popular pet in the world",
  },
  {
    label: "Elephant",
    key: "elephant",
    description: "The largest land animal",
  },
  { label: "Lion", key: "lion", description: "The king of the jungle" },
  { label: "Tiger", key: "tiger", description: "The largest cat species" },
  { label: "Giraffe", key: "giraffe", description: "The tallest land animal" },
  {
    label: "Dolphin",
    key: "dolphin",
    description: "A widely distributed and diverse group of aquatic mammals",
  },
  {
    label: "Penguin",
    key: "penguin",
    description: "A group of aquatic flightless birds",
  },
  {
    label: "Zebra",
    key: "zebra",
    description: "A several species of African equids",
  },
  {
    label: "Shark",
    key: "shark",
    description:
      "A group of elasmobranch fish characterized by a cartilaginous skeleton",
  },
  {
    label: "Whale",
    key: "whale",
    description: "Diverse group of fully aquatic placental marine mammals",
  },
  {
    label: "Otter",
    key: "otter",
    description: "A carnivorous mammal in the subfamily Lutrinae",
  },
  {
    label: "Crocodile",
    key: "crocodile",
    description: "A large semiaquatic reptile",
  },
];

export default function Autocomplete() {
  const [selectKey, setSelectKey] = React.useState<Key | null>(null);

  return (
    <React.Fragment>
      <Input
        aria-hidden
        hidden
        className="invisible hidden"
        name="auto_complete"
        value={selectKey?.toString() || ""}
      />
      <HeroAutocomplete
        className="max-w-xs"
        items={animals}
        label="Favorite Animal"
        placeholder="Search an animal"
        variant="bordered"
        onSelectionChange={setSelectKey}
      >
        {(item) => (
          <AutocompleteItem key={item.key}>{item.label}</AutocompleteItem>
        )}
      </HeroAutocomplete>
    </React.Fragment>
  );
}
