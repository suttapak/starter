"use client";

import { useState } from "react";
import {
  ChevronDown,
  Download,
  Plus,
  Search,
  SlidersHorizontal,
  Trash2,
} from "lucide-react";

import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Badge } from "@/components/ui/badge";

interface Product {
  id: string;
  name: string;
  sku: string;
  price: number;
  status: "In Stock" | "Low Stock" | "Out of Stock";
  category: string;
  rating: number;
  sales: number;
  inventory: number;
  lastUpdated: string;
  vendor: string;
}

const products: Product[] = [
  {
    id: "1",
    name: "Premium Wireless Headphones",
    sku: "WH-001-BLK",
    price: 199.99,
    status: "In Stock",
    category: "Electronics",
    rating: 4.5,
    sales: 1234,
    inventory: 50,
    lastUpdated: "2024-01-20",
    vendor: "AudioTech Pro",
  },
  {
    id: "2",
    name: "Organic Cotton T-Shirt",
    sku: "AT-100-WHT",
    price: 29.99,
    status: "Low Stock",
    category: "Apparel",
    rating: 4.8,
    sales: 2156,
    inventory: 10,
    lastUpdated: "2024-01-22",
    vendor: "EcoWear",
  },
  {
    id: "3",
    name: "Smart Home Security Camera",
    sku: "SC-200-SLV",
    price: 149.99,
    status: "Out of Stock",
    category: "Electronics",
    rating: 4.2,
    sales: 867,
    inventory: 0,
    lastUpdated: "2024-01-21",
    vendor: "SecureLife",
  },
  {
    id: "4",
    name: "Stainless Steel Water Bottle",
    sku: "WB-500-BLU",
    price: 24.99,
    status: "In Stock",
    category: "Accessories",
    rating: 4.7,
    sales: 3421,
    inventory: 75,
    lastUpdated: "2024-01-23",
    vendor: "EcoVessel",
  },
  {
    id: "5",
    name: "Ergonomic Office Chair",
    sku: "OC-300-BLK",
    price: 299.99,
    status: "In Stock",
    category: "Furniture",
    rating: 4.6,
    sales: 543,
    inventory: 25,
    lastUpdated: "2024-01-22",
    vendor: "ErgoComfort",
  },
];

export default function ProductsPage() {
  const [selectedProducts, setSelectedProducts] = useState<Set<string>>(
    new Set(),
  );

  const toggleProduct = (productId: string) => {
    const newSelected = new Set(selectedProducts);

    if (newSelected.has(productId)) {
      newSelected.delete(productId);
    } else {
      newSelected.add(productId);
    }
    setSelectedProducts(newSelected);
  };

  const toggleAll = () => {
    if (selectedProducts.size === products.length) {
      setSelectedProducts(new Set());
    } else {
      setSelectedProducts(new Set(products.map((p) => p.id)));
    }
  };

  const getStatusColor = (status: Product["status"]) => {
    switch (status) {
      case "In Stock":
        return "bg-green-500/10 text-green-500 hover:bg-green-500/20";
      case "Low Stock":
        return "bg-yellow-500/10 text-yellow-500 hover:bg-yellow-500/20";
      case "Out of Stock":
        return "bg-red-500/10 text-red-500 hover:bg-red-500/20";
    }
  };

  return (
    <div className="">
      <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Products</h1>
          <p className="text-muted-foreground">
            Manage and monitor your product inventory
          </p>
        </div>
        <div className="flex items-center gap-4">
          <Button>
            <Plus className="mr-2 h-4 w-4" />
            Add Product
          </Button>
          <Button variant="outline">
            <Download className="mr-2 h-4 w-4" />
            Export
          </Button>
        </div>
      </div>

      <div className="flex flex-col gap-4 md:flex-row md:items-center justify-between mb-6">
        <div className="flex items-center gap-4 flex-1">
          <div className="relative flex-1 md:max-w-sm">
            <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
            <Input
              className="pl-8"
              placeholder="Search products..."
              type="search"
            />
          </div>
          <Button size="icon" variant="outline">
            <SlidersHorizontal className="h-4 w-4" />
          </Button>
        </div>
        <div className="flex items-center gap-4">
          <Select defaultValue="all">
            <SelectTrigger className="w-40">
              <SelectValue placeholder="Category" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All Categories</SelectItem>
              <SelectItem value="electronics">Electronics</SelectItem>
              <SelectItem value="apparel">Apparel</SelectItem>
              <SelectItem value="accessories">Accessories</SelectItem>
              <SelectItem value="furniture">Furniture</SelectItem>
            </SelectContent>
          </Select>
          <Select defaultValue="newest">
            <SelectTrigger className="w-40">
              <SelectValue placeholder="Sort by" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="newest">Newest First</SelectItem>
              <SelectItem value="price-high">Price: High to Low</SelectItem>
              <SelectItem value="price-low">Price: Low to High</SelectItem>
              <SelectItem value="best-selling">Best Selling</SelectItem>
            </SelectContent>
          </Select>
        </div>
      </div>

      <div className="rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="w-12">
                <Checkbox
                  aria-label="Select all products"
                  checked={selectedProducts.size === products.length}
                  onCheckedChange={toggleAll}
                />
              </TableHead>
              <TableHead className="">Product</TableHead>
              <TableHead className="hidden md:table-cell">Category</TableHead>
              <TableHead>Status</TableHead>
              <TableHead className="hidden lg:table-cell">Rating</TableHead>
              <TableHead className="text-right">Price</TableHead>
              <TableHead className="hidden md:table-cell">Inventory</TableHead>
              <TableHead className="hidden lg:table-cell">Sales</TableHead>
              <TableHead className="hidden xl:table-cell">Vendor</TableHead>
              <TableHead className="w-12" />
            </TableRow>
          </TableHeader>
          <TableBody>
            {products.map((product) => (
              <TableRow key={product.id}>
                <TableCell>
                  <Checkbox
                    aria-label={`Select ${product.name}`}
                    checked={selectedProducts.has(product.id)}
                    onCheckedChange={() => toggleProduct(product.id)}
                  />
                </TableCell>
                <TableCell>
                  <div className="flex flex-col">
                    <span className="font-medium">{product.name}</span>
                    <span className="text-sm text-muted-foreground">
                      {product.sku}
                    </span>
                  </div>
                </TableCell>
                <TableCell className="hidden md:table-cell">
                  {product.category}
                </TableCell>
                <TableCell>
                  <Badge
                    className={getStatusColor(product.status)}
                    variant="secondary"
                  >
                    {product.status}
                  </Badge>
                </TableCell>
                <TableCell className="hidden lg:table-cell">
                  <div className="flex items-center">
                    <span className="font-medium">{product.rating}</span>
                    <span className="text-muted-foreground">/5</span>
                  </div>
                </TableCell>
                <TableCell className="text-right">
                  ${product.price.toFixed(2)}
                </TableCell>
                <TableCell className="hidden md:table-cell">
                  {product.inventory.toLocaleString()}
                </TableCell>
                <TableCell className="hidden lg:table-cell">
                  {product.sales.toLocaleString()}
                </TableCell>
                <TableCell className="hidden xl:table-cell">
                  {product.vendor}
                </TableCell>
                <TableCell>
                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button
                        className="h-8 w-8 p-0"
                        size="icon"
                        variant="ghost"
                      >
                        <ChevronDown className="h-4 w-4" />
                        <span className="sr-only">Open menu</span>
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                      <DropdownMenuItem>Edit</DropdownMenuItem>
                      <DropdownMenuItem>Duplicate</DropdownMenuItem>
                      <DropdownMenuSeparator />
                      <DropdownMenuItem className="text-destructive">
                        <Trash2 className="mr-2 h-4 w-4" />
                        Delete
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>
    </div>
  );
}
