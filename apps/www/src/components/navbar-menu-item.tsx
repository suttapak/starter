import { Fragment, useCallback } from "react";
import { Button, cn } from "@heroui/react";
import { Link, useLocation } from "@tanstack/react-router";

export type MNavbarMenuItemsProps = {
  name: string;
  display_name: string;
  path?: string;
  disabled?: boolean;
  external?: boolean;
  icon: React.ReactNode;
  child?: MNavbarMenuItemsProps[];
  padding: number;
  onClose?: undefined | (() => void);
};
const MNavbarMenuItems = ({
  padding,
  disabled,
  display_name,
  path,
  icon,
  child,
  onClose,
}: MNavbarMenuItemsProps) => {
  const { pathname } = useLocation();
  const active = pathname === `${path}`;
  const handleClick = useCallback(() => {
    onClose?.();
  }, []);

  return (
    <Fragment>
      <li className={cn(padding && "pl-10")}>
        <Button
          fullWidth
          as={Link}
          className={"flex justify-start data-[disabled=true]:opacity-100"}
          color={active ? "primary" : "default"}
          isDisabled={disabled || !path}
          size="sm"
          startContent={icon}
          to={path ? path : undefined}
          variant={active ? "solid" : "light"}
          onPress={handleClick}
        >
          {display_name} {disabled && "(disabled)"}
        </Button>
      </li>
      {child && (
        <ul className="flex flex-col gap-1">
          {child.map((item) => (
            <MNavbarMenuItems
              key={item.name}
              child={item.child}
              disabled={item.disabled}
              display_name={item.display_name}
              external={item.external}
              icon={item.icon}
              name={item.name}
              padding={padding ? padding + 4 : 0 + 4}
              path={item.path}
              onClose={onClose}
            />
          ))}
        </ul>
      )}
      {/* {child && <Divider sx={{ backgroundColor: "neutral.400", color: "neutral.400" }} />} */}
    </Fragment>
  );
};

export default MNavbarMenuItems;
