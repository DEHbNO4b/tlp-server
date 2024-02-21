import socket
import select
import json
import time
import datetime
import struct
import sys
import os

path='./Data/LS/'

class TUalf:
    def __init__(self, *args):
        self.Version = int(args[0])
        self.Network = int(args[1])
        self.Year = int(args[2])
        self.Month = int(args[3])
        self.Day = int(args[4])
        self.Hour = int(args[5])
        self.Minute = int(args[6])
        self.Second = int(args[7])
        self.Nanosecond = int(args[8])
        self.Latitude = float(args[9])
        self.Longitude = float(args[10])
        self.Altitude = int(args[11])
        self.Reserved1 = int(args[12])
        self.Peak_current = int(args[13])
        self.Reserved2 = int(args[14])
        self.VHF_Range = int(args[15])
        self.Flash_data = int(args[16])
        self.Number_sensors = int(args[17])
        self.Degrees_freedom = int(args[18])
        self.E_angle = float(args[19])
        self.E_semi_major = float(args[20])
        self.E_semi_minor = float(args[21])
        self.Chi_squared_value = float(args[22])
        self.Risetime = float(args[23])
        self.Peak_zero_time = float(args[24])
        self.Maximum_rate_rise = float(args[25])
        self.Cloud_indicator = int(args[26])
        self.Angle_indicator = int(args[27])
        self.Signal_indicator = int(args[28])
        self.Timing_indicator = int(args[29])

address = "192.168.1.4"
port = 8082
login = b'{ "id": 0, "stream": "4b205c0a-2fd6-4e5e-b4d4-d31eb0e43918" }\r\n'

def tlp_connect():
    while True:
        try:
          # create a TCP socket
          s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
          # now connect to the tlp-serverd process
          s.connect((address, port))

          # Send initial connect message
          s.sendall(login)
          lastSend = time.time()
          # Let's read/write data
          readers = [ s ]
          ok = True

          while ok:
            # Wait up to half a second for data
            ready_to_read, ready_to_write, in_error =  select.select(readers, [ ], readers, 0.5)

            # Read in data and dump it out to console
            for s in ready_to_read:
                line = s.recv(1024).decode("utf-8")
                if line:
                  process_stroke(line)
                else:
                  ok = False

            # Send ID every 10 seconds to keep connection up
            now = time.time()
            if ((now - lastSend) > 10):
                lastSend = now
                s.sendall(login)


          s.close()        
            
        except Exception as e:
          print(f"Error during connection: {e}")
          time.sleep(30)

def process_stroke(stroke):
    if stroke.startswith("9\tKEEP\tALIVE"):
        return
    #2	0	2024	2	14	11	52	34	601633792	39.5925	35.8766	0	0	33	0	0	-1	3	3	22.31	16.30	1.44	0.49	20.9	5.2	4.4	0	1	0	1
    st = stroke.split("\t")
    st[-1] = st[-1].rstrip('\x0D\x0A')
    if(len(st)==30):
      # Создание экземпляра структуры  
      data = TUalf(*st)
      if(data.Flash_data==-1):
        data.Flash_data=255
      # Генерируем данные для записи в файл
      packed_data = struct.pack('=BBHBBBBBifffhhBBBfffffffBBBB', 
        data.Version, data.Network, data.Year, data.Month, data.Day,
        data.Hour, data.Minute, data.Second, data.Nanosecond,
        data.Latitude, data.Longitude, data.Altitude, data.Peak_current,
        data.VHF_Range, data.Flash_data, data.Number_sensors, data.Degrees_freedom, 
        data.E_angle, data.E_semi_major, data.E_semi_minor, data.Chi_squared_value,
        data.Risetime, data.Peak_zero_time, data.Maximum_rate_rise,
        data.Cloud_indicator, data.Angle_indicator, data.Signal_indicator, data.Timing_indicator)
    
      # Получаем текущую дату
#       current_date = datetime.datetime.now()
#       year = str(current_date.year)
#       month = str(current_date.month).zfill(2)  
#       day = str(current_date.day).zfill(2)
#       hour = str(current_date.hourur).zfill(2)
#       minute = str(current_date.minute).zfill(2)
      
      year = str(data.Year)
      month = str(data.Month).zfill(2)
      day = str(data.Day).zfill(2)
      hour = str(data.Hour).zfill(2)
      minute = str(data.Minute).zfill(2)

      # Создаем структуру каталогов
      directory = os.path.join(path+year, month, day, hour)
      os.makedirs(directory, exist_ok=True)

      # Формируем имя файла
      file_name = os.path.join(directory, year+'_'+month+'_'+day+'_'+hour+'_'+minute+'.l80')

      # Записываем данные в файл
      if os.path.exists(file_path):
        # Открываем файл для записи в режиме добавления
        with open(file_name, 'ab') as file:
          file.write(packed_data)
      else:
        # Запись данных в новый бинарный файл
        with open(file_name, 'wb') as file:
          file.write(packed_data)

      print(f'Данные были успешно записаны в файл {file_name}.')
      

if __name__ == "__main__":
    tlp_connect()
