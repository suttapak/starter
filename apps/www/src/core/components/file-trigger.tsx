import { Button, ButtonProps } from "@heroui/react";
import { ForwardedRef, forwardRef, ReactNode, useMemo } from "react";
import { filterDOMProps, useObjectRef } from "@react-aria/utils";

export interface FileTriggerProps
  extends Omit<React.InputHTMLAttributes<HTMLInputElement>, "onSelect"> {
  acceptedFileTypes?: Array<string>;
  allowsMultiple?: boolean;
  defaultCamera?: "user" | "environment";
  invalid?: boolean;
  errorMessage?: string;
  onSelect?: (files: FileList | null) => void;
  buttonProps?: Omit<ButtonProps, "onPress">;
  children?: ReactNode;
}

export const FileTrigger = forwardRef(function FileTrigger(
  props: FileTriggerProps,
  ref: ForwardedRef<HTMLInputElement>,
) {
  let {
    onSelect,
    acceptedFileTypes,
    allowsMultiple,
    defaultCamera,
    buttonProps,
    invalid,
    errorMessage,
    children,
    ...rest
  } = props;
  let inputRef = useObjectRef(ref);
  let domProps = filterDOMProps(rest);

  const message = useMemo(() => {
    if (invalid && errorMessage) {
      return <label className="text-xs text-danger-500">{errorMessage}</label>;
    }

    return null;
  }, [invalid, errorMessage]);

  return (
    <>
      <Button
        {...buttonProps}
        onPress={() => {
          if (inputRef.current?.value) {
            inputRef.current.value = "";
          }
          inputRef.current?.click();
        }}
      >
        {children}
      </Button>
      {message}
      <input
        {...domProps}
        ref={inputRef}
        accept={acceptedFileTypes?.toString()}
        capture={defaultCamera}
        multiple={allowsMultiple}
        style={{ display: "none" }}
        type="file"
        onChange={(e) => onSelect?.(e.target.files)}
      />
    </>
  );
});
