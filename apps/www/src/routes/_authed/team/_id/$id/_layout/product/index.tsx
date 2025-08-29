import { createFileRoute, Link } from "@tanstack/react-router";
import { z } from "zod";
import {
  Pagination,
  TableHeader,
  Table,
  TableColumn,
  TableBody,
  TableRow,
  TableCell,
  Spinner,
  Chip,
  link,
  Input,
  Button,
  button,
  Avatar,
  Card,
  CardBody,
  CardFooter,
  CardHeader,
} from "@heroui/react";
import { PlusCircle, Search } from "lucide-react";
import { useTranslation } from "react-i18next";

import { useGetProducts } from "@/hooks/use-product";
import ProductActionModal from "@/components/product-action-modal";
import { useFormatImageSrc } from "@/hooks/use-format-image-src";

const productSearchParamsSchema = z.object({
  page: z.number().prefault(1),
  limit: z.number().prefault(10),
  code: z.string().optional(),
  name: z.string().optional(),
  uom: z.string().optional(),
});

export const Route = createFileRoute("/_authed/team/_id/$id/_layout/product/")({
  component: RouteComponent,
  validateSearch: productSearchParamsSchema,
  notFoundComponent: NotFound,
});

function NotFound() {
  const { t } = useTranslation();
  const { id } = Route.useParams();

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-50">
      <Card className="max-w-md w-full">
        <CardHeader>{t("not_found.product.title")}</CardHeader>
        <CardBody>
          <p className="mb-4 text-gray-600">
            {t("not_found.product.description")}
          </p>
        </CardBody>
        <CardFooter>
          <Button>
            <Link params={{ id: id }} to="/team/$id/product">
              {t("not_found.product.action")}
            </Link>
          </Button>
        </CardFooter>
      </Card>
    </div>
  );
}
function RouteComponent() {
  const { t } = useTranslation();
  const { id } = Route.useParams();
  const { page, limit } = Route.useSearch();
  const { data, isLoading } = useGetProducts();
  const navigate = Route.useNavigate();
  const { code, name, uom } = Route.useSearch();

  const { format } = useFormatImageSrc();

  return (
    <div className="flex flex-col gap-2">
      <div className="flex justify-end gap-2">
        <Input
          defaultValue={code}
          endContent={<Search />}
          placeholder={t("product.search.state.code.placeholder")}
          onValueChange={(value) =>
            navigate({
              search: {
                code: value,
                page: 1,
              },
            })
          }
        />
        <Input
          defaultValue={name}
          endContent={<Search />}
          placeholder={t("product.search.state.name.placeholder")}
          onValueChange={(value) =>
            navigate({
              search: {
                name: value,
                page: 1,
              },
            })
          }
        />
        <Input
          defaultValue={uom}
          endContent={<Search />}
          placeholder={t("product.search.state.oum.placeholder")}
          onValueChange={(value) =>
            navigate({
              search: {
                uom: value,
                page: 1,
              },
            })
          }
        />
        <Link
          className={button({
            fullWidth: true,
            color: "primary",
          })}
          params={{ id }}
          to="/team/$id/product/new"
        >
          <PlusCircle />
          {t("product.new_product.action")}
        </Link>
      </div>
      <Table
        aria-label="ตารางแสดงหมวดหมู่สินค้า"
        bottomContent={
          data?.data.meta && (
            <div className="flex w-full justify-center gap-2 flex-col sm:flex-row">
              <Pagination
                isCompact
                showControls
                showShadow
                color="secondary"
                page={page || 1}
                total={data?.data.meta.total_page}
                onChange={(page) => navigate({ search: { page } })}
              />
              <label className="flex items-center text-default-400 text-small">
                Rows per page:
                <select
                  className="bg-transparent outline-none text-default-400 text-small"
                  defaultValue={limit}
                  onChange={(e) => {
                    const newLimit = parseInt(e.target.value, 10);

                    navigate({ search: { limit: newLimit, page: 1 } });
                  }}
                >
                  <option value="5">5</option>
                  <option value="10">10</option>
                  <option value="20">20</option>
                  <option value="50">50</option>
                </select>
              </label>
            </div>
          )
        }
      >
        <TableHeader>
          <TableColumn width={80}>{t("product.state.index")}</TableColumn>
          <TableColumn width={80}>{t("product.state.images")}</TableColumn>
          <TableColumn width={160}>{t("product.state.code")}</TableColumn>
          <TableColumn width={200}>{t("product.state.name")}</TableColumn>
          <TableColumn>{t("product.state.description")}</TableColumn>
          <TableColumn width={80}>{t("product.state.price")}</TableColumn>
          <TableColumn width={80}>{t("product.state.oum")}</TableColumn>
          <TableColumn width={300}>{t("product.state.category")}</TableColumn>
          <TableColumn width={80}>Action</TableColumn>
        </TableHeader>
        <TableBody
          items={data?.data.data || []}
          loadingContent={<Spinner />}
          loadingState={isLoading ? "loading" : "idle"}
        >
          {(item) => (
            <TableRow key={item.id}>
              <TableCell>{item.id}</TableCell>
              <TableCell>
                {
                  <Avatar
                    isBordered
                    radius="sm"
                    src={format(item.product_image?.at(-1)?.image.path)}
                  />
                }
              </TableCell>
              <TableCell>
                <Link
                  className={link()}
                  params={{ id, pid: item.id.toString() }}
                  to="/team/$id/product/$pid"
                >
                  {item.code}
                </Link>
              </TableCell>
              <TableCell>{item.name}</TableCell>
              <TableCell>{item.description}</TableCell>
              <TableCell>
                {
                  Intl.NumberFormat("th-TH", {
                    style: "currency",
                    currency: "THB",
                  }).format(item.price / 100) // Convert price from cents to THB
                }
              </TableCell>
              <TableCell>{item.uom}</TableCell>
              <TableCell className="flex flex-wrap gap-0.5">
                {item.product_product_category.map((cate) => (
                  <Chip key={cate.id} color="secondary" variant="bordered">
                    {cate.category.name}
                  </Chip>
                ))}
              </TableCell>

              <TableCell>
                <ProductActionModal item={item} />
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </div>
  );
}
