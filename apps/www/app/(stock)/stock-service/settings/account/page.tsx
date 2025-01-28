"use client";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { CloudUpload, Paperclip } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import {
  FileInput,
  FileUploader,
  FileUploaderContent,
  FileUploaderItem,
} from "@/components/ui/file-upload";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { useGetUserMe } from "@/hooks/use-user";
import { Skeleton } from "@/components/ui/skeleton";

const formSchema = z.object({
  full_name: z.string().min(2),
  username: z.string().min(2),
  email: z.string(),
  image_profile: z.string().min(0),
});

function Page() {
  const { data, isLoading } = useGetUserMe();
  const user = data?.data.data;
  const [files, setFiles] = useState<File[] | null>(null);

  const dropZoneConfig = {
    maxFiles: 5,
    maxSize: 1024 * 1024 * 4,
    multiple: true,
  };
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
  });

  function onSubmit(values: z.infer<typeof formSchema>) {
    try {
    } catch {}
  }
  if (isLoading) {
    return <Skeleton className="h-32 max-w-screen-sm mx-auto" />;
  }

  return (
    <div className="flex justify-center">
      <Card className="w-full max-w-screen-sm">
        <CardHeader>
          <CardTitle>ข้อมูลผู้ใช้งาน</CardTitle>
          <CardDescription>
            รายละเอียดข้อมูลผู้ใช้งาน ชื่อ email
          </CardDescription>
        </CardHeader>
        <CardContent className="grid grid-cols-12 gap-2">
          <div className="col-span-12 sm:col-span-3 items-center justify-start sm:items-start sm:justify-center flex">
            <Avatar className="h-16 w-16 items-center justify-center">
              <AvatarImage src="https://github.com/shadcn.png" />
              <AvatarFallback>CN</AvatarFallback>
            </Avatar>
          </div>
          <div className="col-span-12 sm:col-span-9">
            <Form {...form}>
              <form
                className="gap-2 mx-auto flex flex-col"
                onSubmit={form.handleSubmit(onSubmit)}
              >
                <FormField
                  control={form.control}
                  name="full_name"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>ชื่อเต็ม</FormLabel>
                      <FormControl>
                        <Input
                          readOnly
                          defaultValue={user?.full_name}
                          placeholder="เมธี สุตภักดิ์"
                          type="text"
                          {...field}
                        />
                      </FormControl>
                      <FormDescription>
                        ชื่อ นามสกุล ไม่ต้องใส่คำนำหน้า (นาย นางสาว นาง)
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="username"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Username</FormLabel>
                      <FormControl>
                        <Input
                          readOnly
                          defaultValue={user?.username}
                          placeholder="example"
                          type="text"
                          {...field}
                        />
                      </FormControl>
                      <FormDescription>ชื่อผู้ใช้งาน</FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="email"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Email</FormLabel>
                      <FormControl>
                        <Input
                          readOnly
                          defaultValue={user?.email}
                          placeholder="example@mail.com"
                          type="email"
                          {...field}
                        />
                      </FormControl>
                      <FormDescription>
                        email สำหรับเปลี่ยนรหัสผ่าน
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="image_profile"
                  render={() => (
                    <FormItem>
                      <FormLabel>รูปโปรไฟล์</FormLabel>
                      <FormControl>
                        <FileUploader
                          className="relative bg-background rounded-lg p-2"
                          dropzoneOptions={dropZoneConfig}
                          value={files}
                          onValueChange={setFiles}
                        >
                          <FileInput
                            className="outline-dashed outline-1 outline-slate-500"
                            id="fileInput"
                          >
                            <div className="flex items-center justify-center flex-col p-8 w-full ">
                              <CloudUpload className="text-gray-500 w-10 h-10" />
                              <p className="mb-1 text-sm text-gray-500 dark:text-gray-400">
                                <span className="font-semibold">
                                  Click to upload
                                </span>
                                &nbsp; or drag and drop
                              </p>
                              <p className="text-xs text-gray-500 dark:text-gray-400">
                                SVG, PNG, JPG or GIF
                              </p>
                            </div>
                          </FileInput>
                          <FileUploaderContent>
                            {files &&
                              files.length > 0 &&
                              files.map((file, i) => (
                                <FileUploaderItem key={i} index={i}>
                                  <Paperclip className="h-4 w-4 stroke-current" />
                                  <span>{file.name}</span>
                                </FileUploaderItem>
                              ))}
                          </FileUploaderContent>
                        </FileUploader>
                      </FormControl>

                      <FormMessage />
                    </FormItem>
                  )}
                />
                <Button disabled type="submit">
                  Submit
                </Button>
              </form>
            </Form>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}

export default Page;
