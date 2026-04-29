#import <AudioToolbox/AudioToolbox.h>
#import <Cocoa/Cocoa.h>
#import <CoreAudio/CoreAudio.h>

static AudioDeviceID defaultOutputDevice(void) {
  AudioObjectPropertyAddress addr = {
      kAudioHardwarePropertyDefaultOutputDevice,
      kAudioObjectPropertyScopeGlobal,
      kAudioObjectPropertyElementMain,
  };
  AudioDeviceID deviceID = 0;
  UInt32 size = sizeof(deviceID);
  AudioObjectGetPropertyData(kAudioObjectSystemObject, &addr, 0, NULL, &size,
                             &deviceID);
  return deviceID;
}

static AudioObjectPropertyAddress volumeAddr = {
    kAudioHardwareServiceDeviceProperty_VirtualMainVolume,
    kAudioDevicePropertyScopeOutput,
    kAudioObjectPropertyElementMain,
};

// Returns volume as 0-100. Returns -1 on error.
// https://developer.apple.com/documentation/coreaudio/audioobjectgetpropertydata(_:_:_:_:_:_:)
// https://developer.apple.com/documentation/audiotoolbox/kaudiohardwareservicedeviceproperty_virtualmainvolume
int GetVolume(void) {
  AudioDeviceID device = defaultOutputDevice();

  Float32 vol = 0;
  UInt32 size = sizeof(vol);
  OSStatus status =
      AudioObjectGetPropertyData(device, &volumeAddr, 0, NULL, &size, &vol);
  if (status != noErr)
    return -1;

  return (int)(vol * 100.0f + 0.5f);
}

// Sets volume from 0-100. Returns 0 on success, -1 on error.
// https://developer.apple.com/documentation/coreaudio/audioobjectsetpropertydata(_:_:_:_:_:_:)
int SetVolume(int level) {
  AudioDeviceID device = defaultOutputDevice();

  Float32 vol = (Float32)level / 100.0f;
  OSStatus status = AudioObjectSetPropertyData(device, &volumeAddr, 0, NULL,
                                               sizeof(vol), &vol);
  if (status != noErr)
    return -1;

  return 0;
}

static AudioObjectPropertyAddress muteAddr = {
    kAudioDevicePropertyMute,
    kAudioDevicePropertyScopeOutput,
    kAudioObjectPropertyElementMain,
};

// Returns 1 if muted, 0 if not, -1 on error.
int GetMute(void) {
  AudioDeviceID device = defaultOutputDevice();

  UInt32 muted = 0;
  UInt32 size = sizeof(muted);
  OSStatus status =
      AudioObjectGetPropertyData(device, &muteAddr, 0, NULL, &size, &muted);
  if (status != noErr)
    return -1;

  return (int)muted;
}

// Sets mute state. Returns 0 on success, -1 on error.
int SetMute(int muted) {
  AudioDeviceID device = defaultOutputDevice();

  UInt32 val = (UInt32)muted;
  OSStatus status =
      AudioObjectSetPropertyData(device, &muteAddr, 0, NULL, sizeof(val), &val);
  if (status != noErr)
    return -1;

  return 0;
}

// Simulate a media key press.
// NX_KEYTYPE_PLAY = 16, NX_KEYTYPE_NEXT = 17, NX_KEYTYPE_PREVIOUS = 18
// IOKit/hidsystem/ev_keymap.h
// https://developer.apple.com/documentation/appkit/nsevent/eventtype/systemdefined
// https://developer.apple.com/documentation/coregraphics/1456527-cgeventpost
void MediaKey(int key) {

  // Key down
  NSEvent *event = [NSEvent otherEventWithType:NSEventTypeSystemDefined
                                      location:NSMakePoint(0, 0)
                                 modifierFlags:0xa00
                                     timestamp:0
                                  windowNumber:0
                                       context:nil
                                       subtype:8
                                         data1:(key << 16) | (0x0a << 8)
                                         data2:-1];
  CGEventPost(kCGHIDEventTap, [event CGEvent]);

  // Key up
  event = [NSEvent otherEventWithType:NSEventTypeSystemDefined
                             location:NSMakePoint(0, 0)
                        modifierFlags:0xb00
                            timestamp:0
                         windowNumber:0
                              context:nil
                              subtype:8
                                data1:(key << 16) | (0x0b << 8)
                                data2:-1];
  CGEventPost(kCGHIDEventTap, [event CGEvent]);
}
