//
//  virtualization.h
//
//  Created by codehex.
//

#pragma once

#import <Foundation/Foundation.h>
#import <Virtualization/Virtualization.h>

/* exported from cgo */
void startHandler(void *err, char *id);
void pauseHandler(void *err, char *id);
void resumeHandler(void *err, char *id);
void changeStateOnObserver(int state, char *id);
int listenerShouldAcceptNewConnectionFromSocketDevice(void *listener, void *conn, void *device, char *listenerID);
void connectToPortForSocketDeviceHandler(char *fnID, void *connection, void *error);

@interface Observer : NSObject
- (void)observeValueForKeyPath:(NSString *)keyPath ofObject:(id)object change:(NSDictionary *)change context:(void *)context;
@end

@interface SocketListenerDelegate : NSObject <VZVirtioSocketListenerDelegate>
@property(copy, readwrite) NSString *listenerID;
- (BOOL)listener:(VZVirtioSocketListener *)listener 
shouldAcceptNewConnection:(VZVirtioSocketConnection *)connection 
fromSocketDevice:(VZVirtioSocketDevice *)socketDevice;
@end

/* BootLoader */
void *newVZLinuxBootLoader(const char *kernelPath);
void setCommandLineVZLinuxBootLoader(void *bootLoaderPtr, const char *commandLine);
void setInitialRamdiskURLVZLinuxBootLoader(void *bootLoaderPtr, const char *ramdiskPath);

/* VirtualMachineConfiguration */
bool validateVZVirtualMachineConfiguration(void *config, void **error);
void *newVZVirtualMachineConfiguration(void *bootLoader,
                                    unsigned int CPUCount,
                                    unsigned long long memorySize);
void setEntropyDevicesVZVirtualMachineConfiguration(void *config,
                                                    void *entropyDevices);
void setMemoryBalloonDevicesVZVirtualMachineConfiguration(void *config,
                                                          void *memoryBalloonDevices);
void setNetworkDevicesVZVirtualMachineConfiguration(void *config,
                                                    void *networkDevices);
void setSerialPortsVZVirtualMachineConfiguration(void *config,
                                                 void *serialPorts);
void setSocketDevicesVZVirtualMachineConfiguration(void *config,
                                                   void *socketDevices);
void setStorageDevicesVZVirtualMachineConfiguration(void *config,
                                                    void *storageDevices);
void setDirectorySharingDevicesVZVirtualMachineConfiguration(void *config, void *directorySharingDevices);

/* Configurations */
void *newVZFileHandleSerialPortAttachment(int readFileDescriptor, int writeFileDescriptor);
void *newVZFileSerialPortAttachment(const char *filePath, bool shouldAppend, void **error);
void *newVZVirtioConsoleDeviceSerialPortConfiguration(void *attachment);
void *newVZBridgedNetworkDeviceAttachment(void *networkInterface);
void *newVZNATNetworkDeviceAttachment(void);
void *newVZFileHandleNetworkDeviceAttachment(int fileDescriptor);
void *newVZVirtioNetworkDeviceConfiguration(void *attachment);
void *newVZVirtioEntropyDeviceConfiguration(void);
void *newVZVirtioBlockDeviceConfiguration(void *attachment);
void *newVZDiskImageStorageDeviceAttachment(const char *diskPath, bool readOnly, void **error);
void *newVZVirtioTraditionalMemoryBalloonDeviceConfiguration();
void *newVZVirtioSocketDeviceConfiguration();
void *newVZVirtioFileSystemDeviceConfiguration(const char *tag);

/* VirtualMachine */
void *newVZVirtualMachineWithDispatchQueue(void *config, void *queue, const char *vmid);
bool requestStopVirtualMachine(void *machine, void *queue, void **error);
void startWithCompletionHandler(void *machine, void *queue, const char *vmid);
void pauseWithCompletionHandler(void *machine, void *queue, const char *vmid);
void resumeWithCompletionHandler(void *machine, void *queue, const char *vmid);
bool vmCanStart(void *machine, void *queue);
bool vmCanPause(void *machine, void *queue);
bool vmCanResume(void *machine, void *queue);
bool vmCanRequestStop(void *machine, void *queue);

void *makeDispatchQueue(const char *label);

void *newVZMACAddress(const char *macAddress);
void setNetworkDevicesVZMACAddress(void *config, void *macAddress);
void *getVZBridgedNetworkInterfaces();
char *getVZBridgedNetworkInterfaceLocalizedDisplayName(void *ifPtr);
char *getVZBridgedNetworkInterfaceIdentifier(void *ifPtr);

char *getVZVirtioFileSystemDeviceConfigurationTag(void *ptr);
void *getVZVirtioFileSystemDeviceConfigurationShare(void *ptr);

void *getVZVirtualMachineSocketDevices(void *ptr);
void setSocketListenerForPortVZVirtioSocketDevice(void *ptr, void *listenerPtr, uint32_t port);
void removeSocketListenerForPortVZVirtioSocketDevice(void *ptr, uint32_t port);
void *newVZVirtioSocketListener(const char *listenerID);
uint32_t getVZVirtioSocketConnectionSourcePort(void *ptr);
uint32_t getVZVirtioSocketConnectionDestinationPort(void *ptr);
int getVZVirtioSocketConnectionFileDescriptor(void *ptr);
void closeVZVirtioSocketConnection(void *ptr);
void connectToPortVZVirtioSocketDevice(void *ptr, uint32_t port, const char *fnID);
