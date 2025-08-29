import { CommonModel } from "./common";
import { Image } from "./image_response";
import { ProductCategoryResponse } from "./product_category_response";

export interface ProductCategory extends CommonModel {
  category: ProductCategoryResponse;
}

export interface ProductResponse extends CommonModel {
  team_id: number;
  code: string;
  name: string;
  description: string;
  uom: string;
  price: number;
  product_product_category: ProductCategory[];
  product_image?: ProductImageResponse[];
}

export interface ProductImageResponse extends CommonModel {
  product_id: number;
  image_id: number;
  image: Image;
}
