#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Excel Chiqaradigan Print Server
Version: 2.0
"""

from flask import Flask, request, jsonify, send_file
import openpyxl
from openpyxl.styles import Font, Alignment, Border, Side, PatternFill
import os
from datetime import datetime
import tempfile

app = Flask(__name__)

# Printer sozlamalari
PRINTERS = {
    "p1": "Canon LBP6030",
    "p2": "HP LaserJet 1020", 
    "p3": "Epson L3150",
    "p4": "Brother HL-1110"
}

def create_excel_file(items, username=None, filial=None):
    """Excel fayl yaratish"""
    try:
        # Yangi workbook yaratish
        wb = openpyxl.Workbook()
        ws = wb.active
        ws.title = "Buyurtma"
        
        # Stil sozlamalari
        title_font = Font(name='Arial', size=16, bold=True)
        header_font = Font(name='Arial', size=12, bold=True, color='FFFFFF')
        normal_font = Font(name='Arial', size=11)
        total_font = Font(name='Arial', size=12, bold=True)
        
        # Ranglar
        header_fill = PatternFill(start_color='4472C4', end_color='4472C4', fill_type='solid')
        total_fill = PatternFill(start_color='D9E1F2', end_color='D9E1F2', fill_type='solid')
        
        # Border
        thin_border = Border(
            left=Side(style='thin'),
            right=Side(style='thin'),
            top=Side(style='thin'),
            bottom=Side(style='thin')
        )
        
        # Sarlavha
        ws['A1'] = 'BUYURTMA RO\'YXATI'
        ws['A1'].font = title_font
        ws['A1'].alignment = Alignment(horizontal='center')
        ws.merge_cells('A1:C1')
        
        row = 3
        
        # Ma'lumotlar
        if username:
            ws[f'A{row}'] = 'Buyurtmachi:'
            ws[f'B{row}'] = username
            ws[f'A{row}'].font = Font(bold=True)
            row += 1
        
        if filial:
            ws[f'A{row}'] = 'Filial:'
            ws[f'B{row}'] = filial
            ws[f'A{row}'].font = Font(bold=True)
            row += 1
        
        # Vaqt
        ws[f'A{row}'] = 'Vaqt:'
        ws[f'B{row}'] = datetime.now().strftime('%Y-%m-%d %H:%M:%S')
        ws[f'A{row}'].font = Font(bold=True)
        row += 2
        
        # Jadval sarlavhasi
        ws[f'A{row}'] = '№'
        ws[f'B{row}'] = 'Mahsulot nomi'
        ws[f'C{row}'] = 'Soni'
        
        # Sarlavha stilini qo'llash
        for col in ['A', 'B', 'C']:
            cell = ws[f'{col}{row}']
            cell.font = header_font
            cell.fill = header_fill
            cell.alignment = Alignment(horizontal='center')
            cell.border = thin_border
        
        row += 1
        start_data_row = row
        
        # Mahsulotlar
        total_count = 0
        for i, item in enumerate(items, 1):
            product = item.get('product', 'Noma\'lum')
            count = int(item.get('count', 0))
            total_count += count
            
            ws[f'A{row}'] = i
            ws[f'B{row}'] = product
            ws[f'C{row}'] = count
            
            # Stil qo'llash
            ws[f'A{row}'].alignment = Alignment(horizontal='center')
            ws[f'C{row}'].alignment = Alignment(horizontal='center')
            
            for col in ['A', 'B', 'C']:
                ws[f'{col}{row}'].font = normal_font
                ws[f'{col}{row}'].border = thin_border
            
            row += 1
        
        # Jami qatori
        ws[f'A{row}'] = ''
        ws[f'B{row}'] = 'JAMI:'
        ws[f'C{row}'] = total_count
        
        # Jami qatori stilini qo'llash
        ws[f'B{row}'].font = total_font
        ws[f'C{row}'].font = total_font
        ws[f'B{row}'].alignment = Alignment(horizontal='right')
        ws[f'C{row}'].alignment = Alignment(horizontal='center')
        
        for col in ['A', 'B', 'C']:
            ws[f'{col}{row}'].fill = total_fill
            ws[f'{col}{row}'].border = thin_border
        
        # Ustun kengligini sozlash
        ws.column_dimensions['A'].width = 5
        ws.column_dimensions['B'].width = 40
        ws.column_dimensions['C'].width = 10
        
        # Fayl nomi va yo'li
        timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
        filename = f"buyurtma_{timestamp}.xlsx"
        
        # Temp papkada saqlash
        temp_dir = tempfile.gettempdir()
        filepath = os.path.join(temp_dir, filename)
        
        # Saqlash
        wb.save(filepath)
        
        print(f"✅ Excel fayl yaratildi: {filename}")
        return filepath, filename
        
    except Exception as e:
        print(f"❌ Excel yaratishda xatolik: {e}")
        return None, None

def print_excel_file(filepath, printer_name):
    """Excel faylni chop etish"""
    try:
        import subprocess
        import sys
        
        if sys.platform == "win32":
            # Windows da Excel faylni ochish (default printer ga chop etish)
            os.startfile(filepath, "print")
            print(f"✅ Excel fayl chop etish uchun yuborildi: {printer_name}")
            return True
        else:
            print("❌ Faqat Windows da ishlaydi")
            return False
            
    except Exception as e:
        print(f"❌ Chop etishda xatolik: {e}")
        return False

@app.route('/', methods=['GET'])
def api_info():
    """API ma'lumotlari"""
    return jsonify({
        "name": "Excel Chiqaradigan Print Server",
        "version": "2.0",
        "status": "Ishlamoqda",
        "description": "Excel fayl yaratib chop etadi",
        "printers": PRINTERS,
        "features": [
            "Excel (.xlsx) fayl yaratish",
            "Professional formatlar",
            "Chop etish",
            "Yuklab olish"
        ],
        "endpoints": {
            "POST /print": "Excel yaratib chop etish",
            "POST /excel": "Faqat Excel yaratish (yuklab olish)",
            "GET /": "API ma'lumotlari"
        },
        "example": {
            "url": "/print",
            "method": "POST",
            "body": {
                "printer": "p1",
                "username": "John Doe",
                "filial": "Toshkent",
                "items": [
                    {"product": "Mahsulot 1", "count": 2},
                    {"product": "Mahsulot 2", "count": 1}
                ]
            }
        }
    })

@app.route('/print', methods=['POST'])
def api_print():
    """Excel yaratib chop etish"""
    try:
        data = request.get_json()
        if not data:
            return jsonify({"error": "JSON ma'lumot kerak"}), 400
        
        printer_key = data.get('printer', 'p1').lower()
        items = data.get('items', [])
        username = data.get('username')
        filial = data.get('filial')
        
        if not items:
            return jsonify({"error": "items ro'yxati bo'sh"}), 400
        
        if printer_key not in PRINTERS:
            return jsonify({"error": f"Printer topilmadi: {list(PRINTERS.keys())}"}), 400
        
        printer_name = PRINTERS[printer_key]
        
        # Excel fayl yaratish
        filepath, filename = create_excel_file(items, username, filial)
        if not filepath:
            return jsonify({"error": "Excel yaratib bo'lmadi"}), 500
        
        # Chop etish
        print_success = print_excel_file(filepath, printer_name)
        
        response = {
            "success": True,
            "message": f"Excel fayl yaratildi va {printer_name} ga yuborildi",
            "filename": filename,
            "printer": printer_name,
            "items_count": len(items),
            "total_quantity": sum(int(item.get('count', 0)) for item in items),
            "print_status": "yuborildi" if print_success else "chop etishda xatolik",
            "file_path": filepath
        }
        
        return jsonify(response)
        
    except Exception as e:
        return jsonify({"error": f"Xatolik: {str(e)}"}), 500

@app.route('/excel', methods=['POST'])
def api_excel_only():
    """Faqat Excel yaratish (yuklab olish uchun)"""
    try:
        data = request.get_json()
        if not data:
            return jsonify({"error": "JSON ma'lumot kerak"}), 400
        
        items = data.get('items', [])
        username = data.get('username')
        filial = data.get('filial')
        
        if not items:
            return jsonify({"error": "items ro'yxati bo'sh"}), 400
        
        # Excel fayl yaratish
        filepath, filename = create_excel_file(items, username, filial)
        if not filepath:
            return jsonify({"error": "Excel yaratib bo'lmadi"}), 500
        
        # Faylni yuklab berish
        return send_file(
            filepath,
            as_attachment=True,
            download_name=filename,
            mimetype='application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
        )
        
    except Exception as e:
        return jsonify({"error": f"Xatolik: {str(e)}"}), 500

@app.route('/test', methods=['GET'])
def api_test():
    """Test Excel yaratish"""
    test_items = [
        {"product": "Test mahsulot 1", "count": 2},
        {"product": "Test mahsulot 2", "count": 3},
        {"product": "Test mahsulot 3", "count": 1}
    ]
    
    filepath, filename = create_excel_file(test_items, "Test User", "Test Filial")
    
    if filepath:
        return send_file(
            filepath,
            as_attachment=True,
            download_name=filename,
            mimetype='application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
        )
    else:
        return jsonify({"error": "Test Excel yaratib bo'lmadi"}), 500

if __name__ == '__main__':
    print("=== EXCEL PRINT SERVER ===")
    print("URL: http://localhost:2020")
    print("Test Excel: http://localhost:2020/test")
    print("Printerlar:", list(PRINTERS.keys()))
    print("Excel kutubxonasi: openpyxl")
    print("Ctrl+C - To'xtatish")
    print("===========================")
    
    # openpyxl kutubxonasini tekshirish
    try:
        import openpyxl
        print("✅ openpyxl kutubxonasi mavjud")
    except ImportError:
        print("❌ openpyxl kutubxonasi yo'q!")
        print("💡 O'rnatish: pip install openpyxl")
    
    app.run(host='0.0.0.0', port=2020, debug=False)