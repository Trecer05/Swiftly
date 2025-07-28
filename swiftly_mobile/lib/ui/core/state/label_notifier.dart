import 'dart:ui';

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../domain/models/label_item.dart';
import '../themes/colors.dart';
import 'label_state.dart';

class LabelNotifier extends StateNotifier<LabelState> {
  LabelNotifier() : super(LabelState.initial());

  void addLabel(LabelItem label) {
    state = state.copyWith(labels: [...state.labels, label]);
  }

  void removeLabel(String labelId) {
    state = state.copyWith(
      labels: state.labels.where((label) => label.id != labelId).toList(),
    );
  }

  void updateTitle(String labelId, String newTitle) {
    final processedTitle = newTitle.trim().isEmpty ? 'Название' : newTitle;
    state = state.copyWith(
      labels:
          state.labels.map((label) {
            return label.id == labelId ? label.copyWith(title: processedTitle) : label;
          }).toList(),
    );
  }

  void updateColor(String labelId, Color newColor) {
    state = state.copyWith(
      labels:
          state.labels.map((label) {
            return label.id == labelId ? label.copyWith(color: newColor) : label;
          }).toList(),
    );
  }

  // void loadLabels() async {
  //   state = state.copyWith(isLoading: true);
  //   await Future.delayed(const Duration(seconds: 1));
  //   state = state.copyWith(
  //     labels: [...state.labels, ...mockLabels],
  //     isLoading: false,
  //   );
  // }
}

// final mockLabels = [
//   LabelItem.create(cardId: '3', userId: '1', title: 'чет другое', color: AppColors.yellow),
// ];
